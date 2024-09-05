## Why This Application is Necessary

The DigiTech JamMan XT and JamMan Stereo are powerful loopers, but they face compatibility issues with newer operating systems, especially Windows 11. One major problem is that USB connectivity, which is essential for transferring files between the device and a computer, no longer works reliably on newer systems. This leaves many musicians with limited ways to manage their audio files for these devices.

**This application** was built to solve that problem. Written in Go, it helps users take `.wav` files from their music library and properly format them for transfer to an SSD card that can be read by the JamMan XT or JamMan Stereo. By bypassing USB connectivity and ensuring files are correctly structured, this application offers a reliable way to continue using your DigiTech JamMan on modern operating systems.

But it doesn't stop there.

### Key Features:
- **Automated File Formatting**: This program automatically organizes your `.wav` files into the folder structure that DigiTech JamMan devices require, so there’s no guesswork.
  
- **Customizable Metadata**: You can define details like beats per minute (BPM), rhythm types, and the number of measures, allowing you to customize your backing tracks to suit your performance needs.
  
- **CSV Documentation**: As you import songs, the program also generates a CSV file. This acts as a useful reference for mapping which song goes to which patch location on your JamMan. It not only makes it easy to track your songs but also documents their settings.

### Why This is Different

Most available tools for the JamMan are outdated or don’t work well on new systems. What sets this program apart is that it goes beyond simply copying files. It provides detailed control over your audio tracks, while also offering a convenient way to document and manage your song library with CSV export.

This ensures that even if you’re using older hardware, you can still make the most of your music and backing tracks on modern systems.



## Getting Started

To get this program up and running, follow these steps:

### Prerequisites
Before running the program, make sure you have:
- A MicroSSD card for your DigiTech JamMan XT or JamMan Stereo.
- Your `.wav` files ready for transfer.

### Installation

1. **Format the MicroSD**:
   - Insert the MicroSSD card into your computer.
   - Format it to **FAT32**.
     - **Windows**: Right-click on the drive in File Explorer, select "Format", choose "FAT32" as the file system, and click "Start".
     - **Mac**: Use Disk Utility to erase the MicroSD and set the format to "MS-DOS (FAT)".
   - Name the drive **`JAMMAN`** during the formatting process.

2. **Download the Required Files**:
   - [Download `jamman.exe`](#) and place it in the directory where you want to run the program.
   - [Download `jamman.ini`](#) and configure the settings based on your SSD and file locations.

3. **Configure `jamman.ini`**:
   - Open `jamman.ini` and define the following settings:
     - `JamManType`: Set this to either `JamManSoloXT` or `JamManStereo`, depending on your device.
     - `SSDLocation`: Specify the volume name of your SSD (e.g., `JAMMAN`).
     - `wavFileLoc`: Provide the path where your `.wav` files are stored (e.g., `C:\Users\YourUser\Music\`).

4. **Create the Directory on the MicroSSD**:
   - Ensure that the program will create the correct directory structure (`JamManSoloXT` or `JamManStereo`).

5. **Run the Program**:
   - Open a terminal or command prompt in the directory where `jamman.exe` is located.
   - Run the executable:
     ```bash
     jamman.exe
     ```
   - The program will process your `.wav` files, organize them into the correct folder structure on your SSD, and generate a CSV file documenting the songs, BPMs, rhythm types, and patch locations.

### Compiling from Source

If you prefer to compile the program from source, you can do so using the provided `main.go` file.

1. **Install Go**:
   - Ensure you have [Go installed](https://golang.org/dl/) on your machine.

2. **Clone the Repository**:
   - Clone the repository to your local machine:
     ```bash
     git clone https://github.com/raymondbernard/jammanDigitech.git
     ```

3. **Navigate to the Source Directory**:
   - Change to the directory where `main.go` is located:
     ```bash
     cd jammanDigitech
     ```

4. **Compile the Program**:
   - Use the Go compiler to build the executable:
     ```bash
     go build -o jamman.exe main.go
     ```

5. **Run the Compiled Program**:
   - Now you can run the program by executing the compiled binary:
     ```bash
     jamman.exe
     ```

### Contributions and Suggestions

Contributions and suggestions are always welcome! If you have ideas to improve the program or want to report an issue, feel free to:
- Fork this repository.
- Make your changes.
- Submit a pull request, and we’ll review it.

We’re excited to see how the community can help make this program even better!

