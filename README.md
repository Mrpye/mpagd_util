# mpagd_util

## Overview

`mpagd_util` is a utility tool designed to work with **Multi-Platform Arcade Game Designer (MPAGD)**, a game development tool created by [Jonathan Cauldwell](https://jonathan-cauldwell.itch.io/multi-platform-arcade-game-designer). MPAGD allows developers to create retro-style games for a variety of platforms, including ZX Spectrum, Amstrad CPC, MSX, and others.

This utility provides additional functionality to manage and manipulate MPAGD project files, such as `.apj` files (Arcade Project Files) and `.agd` files (Arcade Game Designer Files). It simplifies tasks like importing, and modifying game assets such as blocks, sprites, screens, and maps.

## Features

- **Read and Write `.apj` Files**: Parse and modify MPAGD project files, including game assets like blocks, sprites, and screens.
- **Import different assets from `.agd` Files**: Import a full `.agd` file or just selected elements into `.apj` files.
- **Auto remapping of sprites and blocks**: Remaps the sprites and blocks if screen is imported and appended.
- **Backup and Restore**: Create backups of project and code files, and restore them when needed to prevent data loss.
- **Auto Backup**: Automatically create backups whenever project files are modified to ensure changes are saved safely.
- **Rotate Sprites and Blocks**: Enables you to easily rotate blocks and sprites and easily create a fully rotated sprite from a single sprite or block.
- **Display Sprites and Blocks**: Render sprites and blocks directly in the terminal or to bitmap for quick visualization.
- **Reorder Sprites and Blocks**: Adjust the sequence of sprites and blocks within a project to better suit your design needs, automatically updating their references throughout the project to maintain consistency.

## Use Cases

- **Game Development**: Enhance your MPAGD projects by importing and managing assets more efficiently.
- **Asset Management**: Reuse or modify blocks, sprites, and other assets across multiple projects.

## Download / Installation

1. Install Go:  
   To use `mpagd_util`, you need to have Go installed on your system. Follow the [official Go installation guide](https://go.dev/doc/install) to set it up.

2. Install the utility:  
   Use the `go install` command to install `mpagd_util` directly:

   ```bash
   go install github.com/Mrpye/mpagd_util.git@latest
   ```

   This will download and install the utility, making it available in your system's PATH.

[You can access the source code here](https://github.com/Mrpye/mpagd_util)

## Documentation

For detailed command-line documentation and examples, refer to the [CLI Documentation](/documents/) for comprehensive guides on all available commands and their usage.

## Limitations

- **Spectrum Projects Only**: Currently, this utility only supports Spectrum-based MPAGD projects. Support for other platforms may be added in the future.
- **Testing Phase**: As the tool is still under testing, ensure you create backups of your project files before using it to avoid accidental data loss.

## How to used

### Auto Backup

Enable auto backup for your project:

```bash
mpagd_util project auto-backup --code
```

### Rotate Sprite

Rotate a sprite 90 degrees clockwise:

```bash
mpagd_util sprites rotate cw [project file] [sprite number] [[output file]]
```

Rotate a sprite 90 degrees counterclockwise:

```bash
mpagd_util sprites rotate ccw [project file] [sprite number] [[output file]]
```

### Import AGD

Import all elements from an AGD file into a project file:

```bash
mpagd_util project import [project file] [agd file] [[output project file]]
```

Import selected elements (e.g., blocks and sprites) from an AGD file:

```bash
mpagd_util project import-selective [project file] [agd file] --blocks --sprites
```

### Import Blocks

Import blocks from an AGD file into a project file:

```bash
mpagd_util blocks import [project file] [agd file] [[output project file]]
```

### Reorder Blocks

Reorder blocks in a project file:

```bash
mpagd_util blocks reorder [project file] "3,1,2,0"
```

## Contributing

We welcome contributions to improve `mpagd_util`. Here's how you can help:

1. **Feature Suggestions**: If you have ideas for new features, feel free to open an issue or submit a pull request.
2. **Code Contributions**: Fork the repository, make your changes, and submit a pull request. Ensure your code follows the project's coding standards.
3. **Documentation**: Help improve the documentation by fixing errors or adding examples.

## Feature Ideas

Here are some ideas for future features:

- Support for additional platforms like Amstrad CPC and MSX.
- Online repo of assets that can be used in project.
- Add offset to reorder so you dont have to enter the index of all the blocks.

## Bug Reporting

If you encounter any issues while using `mpagd_util`, please report them by opening an issue on the [GitHub repository](https://github.com/yourusername/mpagd_util/issues). Include the following details:

- A description of the issue.
- Steps to reproduce the problem.
- Any error messages or logs.
- Your environment (e.g., operating system, MPAGD version).

## Acknowledgments

We would like to thank the following for their contributions and that made it possible for me to make this tool:

- **OSS003**: For sharing the `.apj` file format, which has been instrumental in the development of this utility.
- **Jonathan Cauldwell**: For creating the **Multi-Platform Arcade Game Designer (MPAGD)**, enabling developers to create retro-style games across multiple platforms.
- **The Community**: For their valuable contributions, feedback, and support in improving this tool.
