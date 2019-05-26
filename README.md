# Chrome Native Messaging in Go

Simple Chrome browser extension with a native messaging host written in Go. The focus of this extension and
native messaging host is to showcase the creation of a persistent connection using connectNative() to a native messaging host written in Go and exchange JSON formatted messages.

## Getting Started

The project consists of a Chrome extension app and native messaging host. The native messaging host was written in Go.

### Prerequisites

Chrome v74+

Go v1.10+

Windows 10

### Installing

There are a few steps you must complete to install an unpacked Chrome extension.

**Step 1**: Build the native messaging host exe. Open a terminal and navigate to
the "*native-host/src directory*" in the project. Then, enter the following
command and hit enter:

```
go build -o bin/nativeMsgSample.exe
```

**Step 2**: Update the `/native-host/config/com.sample.native_msg_golang.json` file. Add the full file path of the *nativeMsgSample.exe* file you just created in step 1 to the "path" property value in the JSON file.

Example (change this path to match your file path)...
```
{
    ...
    "path": "C:\\code\\github.com\\chrome-native-messaging-golang\\native-host\\src\\bin\\nativeMsgSample.exe",
    ...
}
```

**Step 3**: Add required registry key to HKCU. Open the Windows Registry Editor (regedit) and navigate to the following path...
```
HKEY_CURRENT_USER/Software/Google/Chrome/NativeMessagingHosts
```
- 3.1: Add a new key with title of `com.sample.native_msg_golang` under the *NativeMessagingHosts* key.
- 3.2: After creating the `com.sample.native_msg_golang` key, there should be a "*(Default)*" string value within the key. Right click on that string value and choose "*Modify*". Then, enter the full path to `/native-host/config/com.sample.native_msg_golang.json`.

**Step 4**: Install the Chrome extension app.

- 4.1: In Chrome, navigate to `chrome://extensions`.
- 4.2: Enable developer mode by toggling the switch in the upper-right corner.
- 4.3: Click on the "Load unpacked" button.
- 4.4: Select the *app* directory in the project to load the html, js, and json files that make up the unpacked extension.

**Step 5**: Run the extension. Open a new tab, and click on the *Apps* button in the Chrome browser toolbar or navigate to `chrome://apps`. Find the "*Chrome Native Messaging Go Example*" app and click on it.

You should see a simple UI containing a button that says "*Connect to Native host*". Click that button to establish a connection to the native messaging host.

Once connected to the native messaging host, a text box and "Send" button should appear in the UI. You can enter "*ping*" into the text box and hit send. This will send a JSON payload containing "*ping*" to the native messaging host. In turn, the host will respond with a JSON payload containing "*pong*".

**Debugging host:** To debug the native messaging host launch Chrome with logging enabled. This will open a terminal window when Chrome is started that may contain messages related to Chrome's interaction with the native messaging host. To enable debugging and view its output, append the `--enable-logging` command to a command to launch chrome, like this: `"C:\Program Files (x86)\Google\Chrome\Application\chrome.exe" --enable-logging`. You can also review the log file the native messaging host will generate. The log file will be found in the same directory as the native messaging host executable.

**Note:** If you do not have a Chrome extension script maintaining a connection to the native messaging host, Chrome will close the Stdin pipe to the host. Depending on how the native messaging host is written, it may or may not close as well. In this sample app, the native host will detect that the Stdin pipe closed and it will trigger the native host to shut down. If the extension is reopened, the native host will start again. I suggest communicating with the native messaging host via a background script. That way, only 1 instance of the native host will be launched.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details