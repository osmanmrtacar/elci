#!/bin/bash

echo "========================================="
echo "TikTok Video Post Test Script"
echo "========================================="

# Step 1: Get JWT token for user ID 1
echo -e "\n1. Getting JWT token for user..."
TOKEN_RESPONSE=$(curl -s http://localhost:8080/api/v1/auth/dev/token/1)
TOKEN=$(echo $TOKEN_RESPONSE | jq -r '.token')

if [ "$TOKEN" == "null" ] || [ -z "$TOKEN" ]; then
    echo "Error: Failed to get token"
    echo $TOKEN_RESPONSE | jq '.'
    exit 1
fi

echo "✓ Got JWT token: ${TOKEN:0:50}..."

# Step 2: Post a video
VIDEO_URL="https://dlp.deplo.xyz/videos/DSlmYsVj-M9.mp4"
CAPTION="Test video posted from API! #test #sosyal"

echo -e "\n2. Posting video to TikTok..."
echo "   Video URL: $VIDEO_URL"
echo "   Caption: $CAPTION"

# Create JSON payload using jq to ensure proper formatting
JSON_PAYLOAD=$(jq -n \
  --arg url "$VIDEO_URL" \
  --arg caption "$CAPTION" \
  '{video_url: $url, caption: $caption}')

POST_RESPONSE=$(curl -s -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "$JSON_PAYLOAD")

echo -e "\nPost Response:"
echo $POST_RESPONSE | jq '.'

# Extract post ID
POST_ID=$(echo $POST_RESPONSE | jq -r '.post.id')

if [ "$POST_ID" == "null" ] || [ -z "$POST_ID" ]; then
    echo -e "\n❌ Failed to create post"
    exit 1
fi

echo -e "\n✓ Post created with ID: $POST_ID"

# Step 3: Check status
echo -e "\n3. Checking post status..."
for i in {1..10}; do
    sleep 3
    STATUS_RESPONSE=$(curl -s http://localhost:8080/api/v1/posts/$POST_ID/status \
      -H "Authorization: Bearer $TOKEN")

    STATUS=$(echo $STATUS_RESPONSE | jq -r '.status')
    echo "   [$i] Status: $STATUS"

    if [ "$STATUS" == "published" ]; then
        echo -e "\n✅ Video published successfully!"
        echo $STATUS_RESPONSE | jq '.'
        break
    elif [ "$STATUS" == "failed" ]; then
        echo -e "\n❌ Video publishing failed"
        echo $STATUS_RESPONSE | jq '.'
        break
    fi
done

echo -e "\n========================================="
echo "Test completed!"
echo "========================================="
