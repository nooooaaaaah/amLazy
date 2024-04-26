Welcome to the amLazy Assistant, the friendly AI-powered TUI (Text-based User Interface) that helps you write shell commands with ease, all while enjoying a virtual bubble tea experience. Ideal for both beginners and seasoned shell users, amLazy offers a unique approach to navigating the complexities of the command line.

## Feature

- **Command Builder**: get a command based on the provided prompt
- **Copy Output**: press crtl-y to copy output and close the app

## TODO

- **Provide ai with some sys info**: Send the os and which shell the user has activated. Makes the prompting less verbose
- **Install via pkg managers**: Create a pkg for brew and others
- **Interactive Command Builder**: Guides you through the process of building shell commands.
- **Refine a comand**: Reprompt so you can add more to a command
- **Educational Mode**: Toggle an educational mode that explains the function and potential uses of different commands and flags.

## Getting Started

To start using amLazy, ensure you have a modern terminal emulator and basic shell access.

1. **Installation**
   As of now just Clone and build. Then copy to bin or wherever your commands are stored.

```sh
git clone https://github.com/nooooaaaaah/amLazy.git
```

Then create a directory amLazy under .config for api keys

```sh
mkdir ~/.config/amLazy
```

With the file config.env

```sh
cd ~/.config/amLazy
touch config.env
```

Then open the confi.env file and add your api key and assistant id

```sh
nvim config.env
```

It should look like this

```.env
OPENAI_API_KEY="{YOUR_API_KEY}"
OPENAI_ASSISTANT_ID="{YOUR_ASSISTANT_ID}"
```

2. **Running amLazy**

Launch amLazy with the following command:

```sh
amLazy
> "Type your question"
```

3. **Using the Assistant**

- **Type your task or command intent**: Such as "compress folder" or "download file from URL".

> Add "on {YOUR SHELL ENV}" for better results

```bash
amLazy
> compress folder to zip
```

- **Execute or Edit**: Once you're satisfied with the suggested command copy it to the clipboard for manual execution.
