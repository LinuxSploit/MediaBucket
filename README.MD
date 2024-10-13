[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/CAR057?referralCode=1136J-)

# MediaBucket Video Transcoding Server

<!-- ![MediaBucket Logo](https://example.com/logo.png) -->

**MediaBucket** is a robust and scalable video transcoding server built with Go. It offers resumable file uploads, automatic video transcoding, and efficient media serving, including thumbnails for both videos and images. Additionally, MediaBucket provides an image upload bucket with automatic thumbnail generation, ensuring seamless media management for your applications.

## Table of Contents

- [Features](#features)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Directory Structure](#directory-structure)
- [Contributing](#contributing)
- [License](#license)
- [Acknowledgements](#acknowledgements)

## Features

- **Resumable File Uploads**: Supports TUS protocol for reliable and resumable uploads of large video and image files.
- **Video Transcoding**: Automatically transcodes uploaded videos into HLS format for adaptive streaming.
- **Thumbnail Generation**:
  - **Videos**: Generates high quality thumbnails and 500px width thumbnails for each uploaded video.
  - **Images**: Generates 500px width thumbnails for each uploaded image.
- **Media Serving**:
  - Serves HLS video streams.
  - Provides access to video and image thumbnails.
  - Hosts an image upload bucket for serving image files.
- **Demo Pages**:
  - `/video-demo`: Demonstrates video upload functionality.
  - `/image-demo`: Demonstrates image upload functionality.
- **CORS Middleware**: Ensures secure and controlled access to media resources.

## Prerequisites

Before setting up MediaBucket, ensure you have the following installed on your system:

- **Go**: Version 1.16 or higher. [Download Go](https://golang.org/dl/)
- **FFmpeg**: Required for video transcoding. [Install FFmpeg](https://ffmpeg.org/download.html)
- **Git**: For cloning the repository. [Download Git](https://git-scm.com/downloads)

## Configuration

MediaBucket uses environment variables for configuration. Below are the essential variables you need to set in your `.env` file:

| Variable          | Description                                     | Example Value             |
| ----------------- | ----------------------------------------------- | ------------------------- |
| `ServerURL`       | Base URL of the server                          | `http://example.com`   |
| `VideoUploadPath` | Endpoint path for video uploads                 | `/video/`                 |
| `ImageUploadPath` | Endpoint path for image uploads                 | `/image/`                 |

**Example `.env` File:**

```env
ServerURL=http://example.com
VideoUploadPath=/video/
ImageUploadPath=/image/
```

## Usage

Once the service is up and running, you can interact with it through the following functionalities:

### 1. Resumable Video Upload

Use the `/video/` endpoint to upload videos. The server supports resumable uploads, allowing you to pause and resume uploads as needed.

### 2. Resumable Image Upload

Use the `/image/` endpoint to upload images. Similar to video uploads, image uploads are resumable.

### 3. Accessing Transcoded Videos

Transcoded HLS video streams are available at the `/hls/` endpoint. For example:

```
http://example.com/hls/{video_id}/master.m3u8
```

### 4. Accessing Thumbnails

- **Video Thumbnails**: Accessible via `/thumbnail/`.
  
  ```
  http://example.com/thumbnail/{video_id}-500w.webp
  ```

- **Image Thumbnails**: Automatically generated with a width of 500px and accessible via the same endpoint.

  ```
  http://example.com/thumbnail/{image_id}-500w.webp
  ```

### 5. Demo Pages

- **Video Upload Demo**: [http://example.com/video-demo](http://example.com/video-demo)
- **Image Upload Demo**: [http://example.com/image-demo](http://example.com/image-demo)

These pages provide a user-friendly interface to test upload functionalities.

## API Endpoints

### Upload Endpoints

- **Video Upload**

  ```
  POST /video/
  ```

  - **Description**: Upload video files using the TUS protocol.
  - **Headers**:
    - `Content-Type: application/offset+octet-stream`
  - **Parameters**: Handled by TUS protocol.

- **Image Upload**

  ```
  POST /image/
  ```

  - **Description**: Upload image files using the TUS protocol.
  - **Headers**:
    - `Content-Type: application/offset+octet-stream`
  - **Parameters**: Handled by TUS protocol.

### Media Access

- **HLS Video Streams**

  ```
  GET /hls/{video_id}/master.m3u8
  ```

  - **Description**: Access the HLS stream of the transcoded video.

- **Thumbnails**

  ```
  GET /thumbnail/{media_id}-500w.webp
  ```

  - **Description**: Retrieve the thumbnail image for the specified media.

### Demo Pages

- **Video Upload Demo**

  ```
  GET /video-demo
  ```

  - **Description**: Serves a web page demonstrating video upload functionality.

- **Image Upload Demo**

  ```
  GET /image-demo
  ```

  - **Description**: Serves a web page demonstrating image upload functionality.

## Directory Structure

```
 /
├── app/
│   ├── main.go              # Main application entry point
│   ├── go.mod               # Go module file
│   └ *                      # rest of the application code
├── storage/
│   ├── MediaBucket/
│   │   ├── videos/          # Original video uploads 
│   │   ├── images/          # Original image uploads
│   │   ├── hls/             # Transcoded HLS streams
│   │   └── thumbnail/       # Generated thumbnails
└── templates/
    ├── video-demo.html      # Video upload demo page
    └── image-demo.html      # Image upload demo page
```

## Contributing

We welcome contributions to enhance MediaBucket! Please follow these steps:

1. **Fork the Repository**

   Click the "Fork" button at the top right of the repository page.

2. **Create a Feature Branch**

   ```bash
   git checkout -b feature/YourFeature
   ```

3. **Commit Your Changes**

   ```bash
   git commit -m "Add your message here"
   ```

4. **Push to Your Fork**

   ```bash
   git push origin feature/YourFeature
   ```

5. **Open a Pull Request**

   Navigate to the original repository and click "Compare & pull request."

Please ensure your code follows the project's coding standards and includes relevant tests.

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgements

- [TUS Protocol](https://tus.io/) for providing a reliable resumable upload protocol.
- [FFmpeg](https://ffmpeg.org/) for powerful multimedia processing capabilities.
- [Go](https://golang.org/) for being the programming language behind MediaBucket.

---

For any questions or support, please open an issue on the [GitHub repository](https://github.com/LinuxSploit/MediaBucket/issues).