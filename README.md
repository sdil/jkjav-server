# JKJAV API Server

This is the implementation of API server described in this [blog post](https://fadhil-blog.dev/blog/how-i-would-built-malaysia-az-site/). TLDR: This API 

## Improvements Can Be Made

- Initialize multiple locations
- Use Redis Lua Scripting to further enhance the atomicity of the API call
- Connect to a Message Queue broker
- Use Go Mutex to ensure the write path is handled synchronously and prevent overbooking
