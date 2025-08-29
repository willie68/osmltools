# osmltools
a small set of tools, working with the data of the openseamap data logger.

mostly written in go, this small tool is for working with the data files from the [open sea map logger](https://wiki.openseamap.org/wiki/OpenSeaMap-dev:HW-logger/OSeaM)

(https://github.com/willie68/OpenSeaMapLogger)

# Install
This is a typically copy/run program. Simply unzip into a folder and you can use it. For easier use, you should save the program to a path that is already entered in the system path, or you can enter the extraction path there.

# Usage
## Check
This will check all files in the sd card folder and write the cleaned files to the output folder. Different errors and warnings will be reported.
Syntax: 
`osml check -d <sd card folder> -f <output folder> [-v] [-o] [-r <report name>]`

-d: folder with the files of the sd card
-f: output folder, where all processed files will be stored
-o: the tool will overwrite existing files
-v: verbose will add more logging output.
-r: the tool will generate a json output file named `report.json` with some additional data

### Processing
First all files of the sd card folder will be parsed, filtered and written to the output folder. Naming of the new files will be
`<vessel id>-<number of file>-<creation date (first GPRMC sentence. in file)>.dat`
After that the tool try to set the creation and change date of the file to the first date found in the nmea sentences. Than , if requested, `osml` will generate the report.
