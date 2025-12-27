#!/bin/bash

echo "========================================="
echo "X (Twitter) Post Test Script"
echo "========================================="

# Check if userId is provided
if [ -z "$1" ]; then
    echo "Usage: ./test_x_post.sh <userId>"
    echo ""
    echo "Steps:"
    echo "1. Login first: open https://adsl-cabinet-acceptance-investigate.trycloudflare.com/api/auth/x/login"
    echo "2. After login, you'll be redirected with userId in URL"
    echo "3. Run this script with your userId: ./test_x_post.sh YOUR_USER_ID"
    exit 1
fi

USER_ID=$1
TEXT="Hello from the X API! Testing automated posting 🚀 #API #automation"

echo "User ID: $USER_ID"
echo "Text: $TEXT"
echo ""

# Create post
echo "Creating post..."
RESPONSE=$(curl -s -X POST http://localhost:8888/api/posts \
  -H "Content-Type: application/json" \
  -d "{
    \"userId\": \"$USER_ID\",
    \"text\": \"$TEXT\"
  }")

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
else
    echo ""
    echo "❌ Failed to create post"
fi

echo ""
echo "========================================="
