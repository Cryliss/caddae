# CADDAE

## Description
This project was inspired by Computer Aided Design and Drafting (CADD) software, as is meant to read *simple* Aerial (AE) redline (hand drawn markings on a construction map) image files and recreate them as a digital running AsBuilt, adding call out boxes that will include information on the work that was performed.

The goal of this program is to see if it's possible to automate the process of creating and maintaining running AsBuilts, which is something I do for my current job as a project coordinator for a construction company.

## Installation
1. Install [Go](https://go.dev/) on your computer, if you haven't already.
2. Download the source code and move into your `go/src/` folder.
3. From terminal, navigate to the project directory, `cd PROJECT_DIRECTORY` and run `make build`.

### Dependencies 
- [gocui](github.com/jroimartin/gocui) `v0.5.0`  
- [errors](github.com/pkg/errors) `v0.9.1`  
- [zerolog](github.com/rs/zerolog) `v1.26.0`  
- [image](golang.org/x/image) `v0.0.0-20210628002857-a66eb6448b8dt`  


Note: This application was built and tested on MacOS.  
Usability is not guaranteed for PC users.

## Starting the Application
The application can be started in one of two ways:
1. Double clicking on the executable file
2. From the project directory in terminal, run `./caddae`

## Instructions for Use
Update each of the following widgets with the requested information

| View Name | Description | Input Format |
| :-------: | :---------- | :----------: |
| Redline | Enter the full path name of the redline file you wish to digitally recreate. | `.png` |
| Running AsBuilt | Enter the full path name of the original asbuilt file that is to be updated. | `.png` |
| DYEA/VZ | Enter the DYEA/VZ# associated with the provided redline. | `DYEA_LSA_8XXXXXX` or `VZ_LAN_0000XXXX` |
| WPD | Enter the date the work was performed. | `MM/DD/YYYY` |
| Production | Enter the quantities for each production unit associated with the redline. | `100` or `100.25` |

Once all information has been entered, click the 'Create!' button to begin the process.

After that, you can view what's happening during the process in the 'Log' panel.

![CADDAE_UI_FILLED_IN](https://github.com/Cryliss/caddae/blob/main/testfiles/Final_UI.png)

## Contribution
Sabra Bilodeau

## License
This project is under the MIT License.

## Project status
Under development.
