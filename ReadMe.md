Welcome to the amLazy Assistant, the friendly AI-powered TUI (Text-based User Interface) that helps you write shell commands with ease. Ideal for both beginners and seasoned shell users, amLazy offers a unique approach to navigating the complexities of the command line.

![](https://raw.githubusercontent.com/nooooaaaaah/amLazy/main/docs/demo.gif)

## Feature

- **Command Builder**: get a command based on the provided prompt
- **Copy Output**: press crtl-y to copy output and close the app

## TODO

- **Better app logging**: need to have togglable logging
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

#### Setup Assistant

Go to open AI's playground to create an assistant.
Feel free to use my instructions:

> you are a tool to generate commands to be run in shell environments like zsh. You can pipe commands together but can not create full scripts. Just give the code for the user to run in their shell. only provide the command in the output.

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
