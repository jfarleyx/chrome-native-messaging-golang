# Chrome Native Messaging in Go

Chrome browser extension with native messaging host written in Go. The focus of this sample extension and
native messaging host is to create a persistent connection to the native messaging host and exchange JSON
with the extension app.

## Getting Started

The project consists of a Chrome extension app and native messaging host. The native messaging host was written in Go.

### Prerequisites

Chrome v74+

Go v1.10+

Windows 10

### Installing

There are a few steps you must complete to install an unpacked Chrome extension.

Step 1: Build the native messaging host exe. Open a terminal and navigate to
the native-host/src directory in the project. Then, enter the following
command and hit enter:

```
go build -o nativeMsgSample.exe
```

Step 2: Update the com.sample.native_msg_golang.json file. Add the full file path to the nativeMsgSample.exe file you just created in step 1 to the "path" property value.

Example...
```
{
    ...
    "path": "C:\\code\\github.com\\chrome-native-messaging-golang\\native-host\\src\\nativeMsgSample.exe",
    ...
}
```

Step 3: Add required registry key to HKCU. Navigate to...
```
HKEY_CURRENT_USER/Software/Google/Chrome/NativeMessagingHosts
```
Add a new key with title of "com.sample.native_msg_golang" under the *NativeMessagingHosts* key.

After creating the `com.sample.native_msg_golang` key, there should be a "(Default)" string value under the key. Right click on that string value and choose Modify. Then, enter the full path to the `com.sample.native_msg_golang.json` file found in the project under the *native-host/config* directory.

Step 4: Install the Chrome extension app.

- 4.1: In Chrome, navigate to `chrome://extensions`.
- 4.2: Enable developer mode by toggling the switch in the upper-right corner.
- 4.3: Click on the "Load unpacked" button.
- 4.4: Select the *extension* directory in the project to load the html, js, and json files that make up the unpacked extension.

Step 5: Run the extension. Open a new tab, and click on the *Apps* button in the Chrome browser toolbar or navigate to `chrome://apps`. Find the "*Chrome Native Messaging Go Example*" app and click on it.

You should see a simple UI containing a button that says "*Connect to Native host*". Click that button to establish a connection to the native messaging host.

Once connected to the native messaging host, a text box and Send button should appear in the UI. You can enter "*ping*" into the text box and hit send. You should see the following...

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details