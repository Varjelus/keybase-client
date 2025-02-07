## Keybase

Hi, and welcome to the Keybase client repo.  All our client apps (OSX, Windows, Linux, iOS, and Android) are being actively developed in this repository. Please, dig around.

### Warnings

We'd love you to read our source code.

But - some of the things in this repo are explorations, and the app you build from source just *might not do what it says it's doing*. So, if you just want to install Keybase on your computer, you should **[monitor our releases](https://keybase.io/download)** for OSX, Linux, or Windows.

![Sharing](https://keybase.io/images/github/repo_share.png?)


### Code Layout

* **go**: Core crypto libraries; the Keybase service; the command line client. [Learn More](go/README.md)
* **react-native**: Android and iOS apps developed with [React Native](https://facebook.github.io/react-native/).
* **desktop**: Desktop application for OSX, Linux and Windows, made with the [Electron](https://github.com/atom/electron) framework, sharing React code with react-native.
* **packaging**: Scripts for releasing packages across the various platforms.
* **protocol**: Defines the protocol for communication for clients to the Keybase services. Uses [Avro](http://avro.apache.org/docs/1.7.7/). [Learn More](protocol/README.md)
* **media**: Icons, graphics, media for Keybase apps.
* **osx**: The Mac OS X Keybase.app, development parallel to an Electron-based application above. [Learn More](osx/README.md)


### Problems?

Report any issues with client software on this GitHub [issue tracker](https://github.com/keybase/client/issues).
Internally, we track our progress using Jira, but all PRs come through GitHub for your review!

If you're having problem with our Website, try the
[keybase-issues](https://github.com/keybase/keybase-issues) issue tracker.

We check and update both frequently.

### License

Most code is released under the New BSD (3 Clause) License.  If subdirectories include
a different license, that license applies instead.

### Cryptography Notice

This distribution includes cryptographic software. The country in which you currently reside may have restrictions on the import, possession, use, and/or re-export to another country, of encryption software. BEFORE using any encryption software, please check your country's laws, regulations and policies concerning the import, possession, or use, and re-export of encryption software, to see if this is permitted. See http://www.wassenaar.org/ for more information.

The U.S. Government Department of Commerce, Bureau of Industry and Security (BIS), has classified this software as Export Commodity Control Number (ECCN) 5D002.C.1, which includes information security software using or performing cryptographic functions with asymmetric algorithms. The form and manner of this distribution makes it eligible for export under the License Exception ENC Technology Software Unrestricted (TSU) exception (see the BIS Export Administration Regulations, Section 740.13) for both object code and source code.
