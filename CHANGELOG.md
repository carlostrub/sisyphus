# Release 1.2.0
## Added
- SISYPHUS_DRY_RUN flag to allow dry runs without moving files. In
  fact, it does create missing folders, the database, learns, and
  classifies mails. But, sisyphus does not move files between
  folders. (#8)

## Changed
- 

## Fixed
- Do not require config to get to help (#7)

## Known Issues
- There seems to be an issue with quotedprintable not properly reading in
  malformed mails. Currently, such is likely to pass the filter.


# Release 1.1.1
## Added
-

## Changed
- 

## Fixed
- Fixed the build process (#5)
- Fixed the dependency on urfave/cli (#5)

## Known Issues
- There seems to be an issue with quotedprintable not properly reading in
  malformed mails. Currently, such is likely to pass the filter.

# Release 1.1.0
## Added
- Perform a database backup before starting a new learning cycle
- Provide a 'stats' command to display various statistics in an info log

## Changed
- 

## Fixed
- 

## Known Issues
- There seems to be an issue with quotedprintable not properly reading in
  malformed mails. Currently, such is likely to pass the filter.

# Release 1.0.0
## Added
-

## Changed
- update dependencies
- application is stable enough to be released as version 1.0.0

## Fixed
- 

## Known Issues
- There seems to be an issue with quotedprintable not properly reading in
  malformed mails. Currently, such is likely to pass the filter.

# Release 0.3.0
## Added
-

## Changed
- Converted the entire app to a [Twelve-Factor App](https://12factor.net/).
  This has consequences in how you launch it, i.e. use environment variables
  instead of flags.
- The interval between learning periods can be set at runtime now.
- Unload mail content after classification and learning, should reduce memory
  requirements.

## Fixed
- Only permit unicode characters of bitsize larger than 2, this guarantees we
  are only accepting for example Chinese characters as individual words. The
  unicode parser introduced in Version 0.2.0 led to individual accented
  characters being falsely treated as a word.

## Known Issues
- There seems to be an issue with quotedprintable not properly reading in
  malformed mails. Currently, such is likely to pass the filter.

# Release 0.2.0
## Added
Support for unicode characters.

## Changed
-

## Fixed
-

## Known Issues
-

# Release 0.1.0
## Added
First working release. Let's fight junk mail!!! Have fun ;-)

## Changed
-

## Fixed
-

## Known Issues
-
