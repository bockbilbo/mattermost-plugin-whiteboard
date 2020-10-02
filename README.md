# Mattermost Whiteboard Plugin

Share virtual whiteboards on [Mattermost](https://mattermost.org/) to collaborate with your team.

## Supported Whiteboard Tools
This plugin works with both public and self-hosted instances of the following open-source whiteboard tools:

- [cracker0dks/whiteboard](https://github.com/cracker0dks/whiteboard)
- [Excalidraw](https://github.com/excalidraw/excalidraw)
- [WBO](https://github.com/lovasoa/whitebophir)

## Screenshot
![image](https://user-images.githubusercontent.com/13318033/94899210-c0470100-0447-11eb-92ae-c9af86b98e87.png)

## How to use the plugin
- Run the `/whiteboard` command to share a virtual whiteboard
    * You can optionally set the whiteboard name appending an ID after the command
    * Custom whiteboard names are not available when using Excalidraw
- Click the whiteboard icon in a channel header to share a new whiteboard

## Installation

1. Access the [releases page](https://github.com/bockbilbo/mattermost-plugin-whiteboard/releases) and download the latest tar.gz file.
2. Upload the file to your Mattermost server through **System Console > Plugins > Management**, or by manually uploading the uncompressed file to the plugin directory of the server.

Read the official [Mattermost documentation](https://docs.mattermost.com/administration/plugins.html#set-up-guide) for more details.

## Configuration

Log into Mattermost with an admin account, go to **System Console > Plugins > Whiteboard** and set the following values:

1. **Enable Plugin**: ``true``
2. **Whiteboard Tool**: select your preferred whiteboard backend from the list
3. **Custom URL**: if using a self-hosted tool, enter the URL to the resource.
    * *Note that the Whiteboard Tool selected on step #2 must match the tool installed in your server*

## Development

This plugin has been a personal project to get familiar with Go. It has been put together with snippets from the [Mattermost plugin starter template](https://github.com/mattermost/mattermost-plugin-starter-template), the [original Mattermost Jitsi Plugin](https://github.com/appmodule/mattermost-plugin-jitsi), the [official Mattermost Jitsi Plugin](https://github.com/mattermost/mattermost-plugin-jitsi), and the [Mattermost Zoom Plugin](https://github.com/mattermost/mattermost-plugin-zoom).

For those of you looking for [Mattermost plugin coding examples](https://developers.mattermost.com/extend/plugins/example-plugins/), this plugin contains both a server and web app portion, and demonstrates the following topics:

* Defining a settings schema, allowing system administrators to configure the plugin via system console UI.
* Extending existing webapp components to add elements to the UI.
* Creating rich posts using custom post types.
* Registering a custom slash command..
* Creating a bot user to post notices in the channels.
* Using a custom HTTP handler to generate and serve content.

### Compiling

1. Set your development environment for [Go](https://golang.org) and [npm](https://www.npmjs.com)
2. Download the source code:
  `git clone https://github.com/bockbilbo/mattermost-plugin-whiteboard.git`
3. Run `make` to compile while checking the quality of the code using [go linterns](https://medium.com/wesionary-team/introduction-to-linters-in-go-and-know-about-golangci-lint-6486fb2864b3)
