#!/bin/bash

echo "========================================="
echo "X (Twitter) Post with Media Test Script"
echo "========================================="

# Check if userId is provided
if [ -z "$1" ]; then
    echo "Usage: ./test_x_post_with_media.sh <userId> [media_url]"
    echo ""
    echo "Examples:"
    echo "  Text only:"
    echo "    ./test_x_post_with_media.sh YOUR_USER_ID"
    echo ""
    echo "  With image:"
    echo "    ./test_x_post_with_media.sh YOUR_USER_ID https://example.com/image.jpg"
    echo ""
    echo "  With video:"
    echo "    ./test_x_post_with_media.sh YOUR_USER_ID https://example.com/video.mp4"
    echo ""
    echo "Steps to get userId:"
    echo "1. Login first: open https://adsl-cabinet-acceptance-investigate.trycloudflare.com/api/auth/x/login"
    echo "2. After login, check URL or run: curl http://localhost:8888/api/auth/x/users"
    exit 1
fi

USER_ID=$1
MEDIA_URL=$2

if [ -z "$MEDIA_URL" ]; then
    # Text only post
    TEXT="Hello from X API! Testing text-only posting 🚀 #API #automation"
    echo "User ID: $USER_ID"
    echo "Text: $TEXT"
    echo "Media: None"
    echo ""

    RESPONSE=$(curl -s -X POST http://localhost:8888/api/posts \
      -H "Content-Type: application/json" \
      -d "{
        \"userId\": \"$USER_ID\",
        \"text\": \"$TEXT\"
      }")
else
    # Post with media
    TEXT="Testing X API with media! 📸🎥 #API #automation"
    echo "User ID: $USER_ID"
    echo "Text: $TEXT"
    echo "Media: $MEDIA_URL"
    echo ""

    RESPONSE=$(curl -s -X POST http://localhost:8888/api/posts \
      -H "Content-Type: application/json" \
      -d "{
        \"userId\": \"$USER_ID\",
        \"text\": \"$TEXT\",
        \"media_urls\": [\"$MEDIA_URL\"]
      }")
fi

echo "Creating post..."
echo ""
echo "Response:"
echo $RESPONSE | jq '.'

# Check if successful
if echo $RESPONSE | jq -e '.success' > /dev/null 2>&1; then
    POST_ID=$(echo $RESPONSE | jq -r '.post.id')
    echo ""
    echo "✅ Post created successfully!"
    echo "Tweet ID: $POST_ID"
    echo "View at: https://twitter.com/user/status/$POST_ID"

    if [ -n "$MEDIA_URL" ]; then
        MEDIA_IDS=$(echo $RESPONSE | jq -r '.media_ids[]')
        echo "Media IDs: $MEDIA_IDS"
    fi
else
    echo ""
    echo "❌ Failed to create post"
fi

echo ""
echo "========================================="
