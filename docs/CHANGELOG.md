# Changelog

All notable changes to this project (hashcat.launcher) will be documented in this file.

## [dev] DD/MM/YYYY
### Added:
 - i18n support (#77)
 - a dropdown with latest used masks (#43)

### Fixed:
 - Queue is not taking the first task added (#84)

## [1.1.1] 06/07/2022
### Changed:
 - Dialog handler

### Fixed:
 - remove tmp folder on exit (#89)
 - duplicated info (session name) (#81)
 - cap to hcwpax tool's output name is only based on one wpa even if there is many (#78)
 - removed 'Chrome is being controlled by automated test software' message (#92)

## [1.1.0] 27/12/2021
### Added:
 - Directory watcher (to auto load new added files) (#67)
 - Option to save cap to hcwpax output to hashes
 - An alert if hashcat is not found

### Changed:
 - Moved default extra arguments to general (#66)
 - Aborted status color to yellow (#71)
 - Finished tasks can be started again
 - Task ID is now based on an incremental number

### Fixed:
 - Unable to resume new tasks that just hit checkpont (#74)

## [1.0.0] 22/12/2021
### Added:
 - Missing Started and ETA info for tasks
 - Task guess info (guess base and guess mod) (#53)

### Changed:
 - Disable potfile by default

## [1.0.0-beta] 20/12/2021
### Changed:
 - Exhausted status color to pink (#58)
 - Export configs to /exported directory (#59)

### Added:
 - Allow loading files from symlinks
 - Allow preserving task config
 - Import/Export task config feature
 - Auto load existing .restore files
 - Queued tasks system
 - Delete task from tasks list

### Fixed:
 - Missing hybrid attacks
 - Mask file input not displaying the correct set value
 - Arguments list not scrollable (#56)

## [1.0.0-alpha] 12/12/2021
### Changed:
 - refactored the code
 - design

### Added:
 - Filter/Select devices by id (#40)
 - Tools section

## [0.5.2] 28/06/2021
### Changed:
 - hashcat min required version: v6.2.1

### Added:
 - Markov options (#31)
 - Two buttons to add & select dictionaries in one click

### Fixed:
 - A versioning typo

## [0.5.1] 13/03/2021
### Added:
 - User confirmation before resetting stats
 - Delete restore files feature
 - More scaling options

### Fixed:
 - Arguments list in the restore modal was only read line by line
 - UI too big and scaling issues

## [0.5.0] - 10/03/2021
### Changed:
 - UI Design

## [0.4.0] - 21/12/2020
### Changed:
 - Search for hashtypes by hashcat mode (name and id)
 - Tasks Journals now hold all the info

### Added:
 - Notifications Feature
 - Priority Feature
 - Restore Feature
 - Get Info about a task
 - Skip/Bypass an attack

## [0.3.1] - 21/11/2020
### Fixed:
 - File Dialog not working for some linux users

### Changed:
 - File Dialog
 - Themes and Sizes

## [0.3.0] - 17/08/2020
### Changed:
 - Enhanced files selection for dictionaries and rules

### Added:
 - FileBase Feature

### Fixed:
 - Hardware Monitoring: Wrong hardware id
 - Hardware Monitoring: Missing some stats

## [0.2.0] - 17/07/2020
### Changed:
 - hashcat min required version: v6.0.0

### Added:
 - Can select a whole folder of dictionaries
 - Added few more extensions to the dictionary selection dialog's filter
 - Hardware Monitoring

## [0.1.2] - 25/04/2020
### Fixed:
 - Spelling errors
 - Fields were limited and can't type long arguments

### Changed:
 - Hash Type search field is now fixed at the top
 - Merged appearanceScreen to optionsScreen

### Added:
 - Mask attack supports hcmask files
 - Mask length counter
 - Help messages

## [0.1.1] - 21/04/2020
### Fixed:
 - Tabs texts were not being updated
 - Session ID was not being set correctly
 - Sessions were not being removed on tab close
 - Windows: Incorrect working directory
 - Windows: Broken control functions

## [0.1.0] - 17/04/2020
- Initial Release
___
Date format: DD/MM/YYYY