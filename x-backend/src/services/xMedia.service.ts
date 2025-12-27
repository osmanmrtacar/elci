import axios from 'axios'
import FormData from 'form-data'
import fs from 'fs'
import { promisify } from 'util'

const sleep = promisify(setTimeout)

interface InitResponse {
  data: {
    id: string
    media_key: string
    expires_after_secs: number
    processing_info?: {
      state: 'pending' | 'in_progress' | 'failed' | 'succeeded'
      check_after_secs?: number
      progress_percent?: number
    }
  }
}

interface FinalizeResponse {
  data: {
    id: string
    media_key: string
    expires_after_secs: number
    processing_info?: {
      state: 'pending' | 'in_progress' | 'failed' | 'succeeded'
      check_after_secs?: number
      progress_percent?: number
    }
    size: number
  }
}

interface StatusResponse {
  data: {
    id: string
    processing_info?: {
      state: 'pending' | 'in_progress' | 'failed' | 'succeeded'
      check_after_secs?: number
      progress_percent?: number
    }
  }
}

export class XMediaService {
  private readonly uploadUrl = 'https://api.x.com/2/media/upload'
  private readonly chunkSize = 512 * 1024 // 512KB chunks (X API limit)

  // Download file from URL to temporary location
  async downloadFile(url: string, outputPath: string): Promise<void> {
    const response = await axios({
      method: 'GET',
      url: url,
      responseType: 'stream',
    })

    const writer = fs.createWriteStream(outputPath)
    response.data.pipe(writer)

    return new Promise((resolve, reject) => {
      writer.on('finish', resolve)
      writer.on('error', reject)
    })
  }

  // Get media type from file extension
  private getMediaType(filePath: string): string {
    const ext = filePath.toLowerCase().split('.').pop()
    const typeMap: Record<string, string> = {
      jpg: 'image/jpeg',
      jpeg: 'image/jpeg',
      png: 'image/png',
      gif: 'image/gif',
      webp: 'image/webp',
      bmp: 'image/bmp',
      tiff: 'image/tiff',
      mp4: 'video/mp4',
      mov: 'video/quicktime',
      webm: 'video/webm',
      ts: 'video/mp2t',
      srt: 'text/srt',
      vtt: 'text/vtt',
    }
    // Default to video/mp4 if extension not found
    return typeMap[ext || ''] || 'video/mp4'
  }

  // Get media category based on type
  private getMediaCategory(mediaType: string): string {
    if (mediaType.startsWith('video/')) {
      return 'amplify_video'
    }
    if (mediaType.startsWith('text/')) {
      return 'subtitles'
    }
    if (mediaType === 'image/gif') {
      return 'tweet_gif'
    }
    return 'tweet_image'
  }

  // Step 1: Initialize chunked upload (v2 API with OAuth 2.0)
  async initUpload(
    accessToken: string,
    filePath: string
  ): Promise<InitResponse> {
    const stats = fs.statSync(filePath)
    const mediaType = this.getMediaType(filePath)
    const mediaCategory = this.getMediaCategory(mediaType)

    const requestBody = {
      media_type: mediaType,
      total_bytes: stats.size,
      media_category: mediaCategory,
    }

    console.log('Initializing upload (dedicated endpoint /initialize):', requestBody)

    try {
      const response = await axios.post(
        `${this.uploadUrl}/initialize`,
        requestBody,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
          },
        }
      )

      console.log('Upload initialized:', response.data.data.id)
      return response.data
    } catch (error: any) {
      console.error('Init upload error:', error.response?.data || error.message)
      throw new Error(`Failed to initialize upload: ${error.message}`)
    }
  }

  // Step 2: Upload file chunks (v2 API with OAuth 2.0)
  async appendChunks(
    accessToken: string,
    mediaId: string,
    filePath: string
  ): Promise<void> {
    const stats = fs.statSync(filePath)
    const totalChunks = Math.ceil(stats.size / this.chunkSize)

    // Read entire file into buffer once
    const fileBuffer = fs.readFileSync(filePath)
    console.log(`File loaded: ${fileBuffer.length} bytes`)
    console.log(`Uploading ${totalChunks} chunks...`)

    for (let segmentIndex = 0; segmentIndex < totalChunks; segmentIndex++) {
      const start = segmentIndex * this.chunkSize
      const end = Math.min(start + this.chunkSize, stats.size)

      // Slice the buffer to get the chunk
      const chunk = fileBuffer.slice(start, end)

      console.log(`Chunk ${segmentIndex}: ${chunk.length} bytes (${start}-${end})`)

      const formData = new FormData()
      formData.append('media', chunk, {
        filename: 'blob',
        contentType: 'application/octet-stream',
      })
      formData.append('segment_index', segmentIndex.toString())

      try {
        await axios.post(`${this.uploadUrl}/${mediaId}/append`, formData, {
          headers: {
            ...formData.getHeaders(),
            Authorization: `Bearer ${accessToken}`,
          },
          maxContentLength: Infinity,
          maxBodyLength: Infinity,
        })

        console.log(
          `Uploaded chunk ${segmentIndex + 1}/${totalChunks} (${Math.round((end / stats.size) * 100)}%)`
        )
      } catch (error: any) {
        console.error(`Append chunk ${segmentIndex} error:`)
        console.error('Status:', error.response?.status)
        console.error('Headers:', error.response?.headers)
        console.error('Data:', JSON.stringify(error.response?.data, null, 2))
        console.error('Message:', error.message)
        throw new Error(`Failed to upload chunk ${segmentIndex}: ${error.message}`)
      }
    }
  }

  // Step 3: Finalize upload (v2 API with OAuth 2.0)
  async finalizeUpload(
    accessToken: string,
    mediaId: string
  ): Promise<FinalizeResponse> {
    try {
      const response = await axios.post(
        `${this.uploadUrl}/${mediaId}/finalize`,
        {},
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
            'Content-Type': 'application/json',
          },
        }
      )

      console.log('Upload finalized:', response.data)
      return response.data
    } catch (error: any) {
      console.error(
        'Finalize upload error:',
        error.response?.data || error.message
      )
      throw new Error(`Failed to finalize upload: ${error.message}`)
    }
  }

  // Step 4: Check processing status (for videos)
  async checkStatus(
    accessToken: string,
    mediaId: string
  ): Promise<StatusResponse> {
    try {
      const response = await axios.get(this.uploadUrl, {
        params: {
          command: 'STATUS',
          media_id: mediaId,
        },
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      })

      return response.data
    } catch (error: any) {
      console.error('Status check error:', error.response?.data || error.message)
      throw new Error(`Failed to check status: ${error.message}`)
    }
  }

  // Wait for processing to complete
  async waitForProcessing(
    accessToken: string,
    mediaId: string,
    maxWaitSeconds: number = 300
  ): Promise<void> {
    const startTime = Date.now()

    while (true) {
      const status = await this.checkStatus(accessToken, mediaId)

      if (!status.data.processing_info) {
        // No processing needed
        return
      }

      const state = status.data.processing_info.state

      if (state === 'succeeded') {
        console.log('Media processing completed')
        return
      }

      if (state === 'failed') {
        throw new Error('Media processing failed')
      }

      // Check timeout
      const elapsedSeconds = (Date.now() - startTime) / 1000
      if (elapsedSeconds > maxWaitSeconds) {
        throw new Error('Media processing timeout')
      }

      // Wait before next check
      const waitSeconds = status.data.processing_info.check_after_secs || 5
      console.log(
        `Processing ${state} (${status.data.processing_info.progress_percent || 0}%), checking again in ${waitSeconds}s...`
      )
      await sleep(waitSeconds * 1000)
    }
  }

  // Complete upload process: download, upload, and return media_id
  async uploadFromUrl(accessToken: string, mediaUrl: string): Promise<string> {
    const tempFile = `/tmp/x-media-${Date.now()}-${Math.random().toString(36).substring(7)}`

    try {
      // Download file
      console.log('Downloading media from:', mediaUrl)
      await this.downloadFile(mediaUrl, tempFile)

      // Initialize upload (v2 API)
      const initResponse = await this.initUpload(accessToken, tempFile)
      const mediaId = initResponse.data.id

      // Upload chunks (v2 API)
      await this.appendChunks(accessToken, mediaId, tempFile)

      // Finalize upload (v2 API)
      const finalizeResponse = await this.finalizeUpload(accessToken, mediaId)

      // Wait for processing if needed
      if (finalizeResponse.data.processing_info) {
        await this.waitForProcessing(accessToken, mediaId)
      }

      console.log('Media uploaded successfully:', mediaId)
      return mediaId
    } finally {
      // Clean up temp file
      if (fs.existsSync(tempFile)) {
        fs.unlinkSync(tempFile)
      }
    }
  }
}
