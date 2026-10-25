; Snugfit installer (NSIS 3, Modern UI 2)
; Build:  makensis -DVERSION=1.0.0 -DSRCDIR=..\dist\Snugfit installer\snugfit.nsi
; Produces Snugfit-Setup-<version>.exe next to this script (or OUTDIR if given).

Unicode true
ManifestDPIAware true
SetCompressor /SOLID lzma
RequestExecutionLevel user

!ifndef VERSION
  !define VERSION "1.0.0"
!endif
!ifndef SRCDIR
  !define SRCDIR "..\dist\Snugfit"
!endif
!ifndef OUTDIR
  !define OUTDIR "."
!endif

!define APPNAME "Snugfit"
!define PUBLISHER "Adam Glassner"
!define WEBSITE "https://github.com/adamglassner/snugfit"
!define UNINSTKEY "Software\Microsoft\Windows\CurrentVersion\Uninstall\Snugfit"

Name "${APPNAME}"
OutFile "${OUTDIR}\Snugfit-Setup-${VERSION}.exe"
InstallDir "$LOCALAPPDATA\Programs\Snugfit"
InstallDirRegKey HKCU "Software\Snugfit" "InstallDir"
BrandingText "Snugfit ${VERSION}  -  Relaxing jigsaw puzzles"
ShowInstDetails nevershow
ShowUninstDetails nevershow

VIProductVersion "${VERSION}.0"
VIAddVersionKey "ProductName" "${APPNAME}"
VIAddVersionKey "CompanyName" "${PUBLISHER}"
VIAddVersionKey "FileDescription" "Snugfit Setup"
VIAddVersionKey "FileVersion" "${VERSION}"
VIAddVersionKey "ProductVersion" "${VERSION}"
VIAddVersionKey "LegalCopyright" "(c) 2026 ${PUBLISHER}"

!include "MUI2.nsh"
!include "FileFunc.nsh"

; ---------------------------------------------------------------- theme
!define MUI_ICON "..\branding\snugfit.ico"
!define MUI_UNICON "..\branding\snugfit.ico"
!define MUI_HEADERIMAGE
!define MUI_HEADERIMAGE_BITMAP "..\branding\installer_header.bmp"
!define MUI_HEADERIMAGE_UNBITMAP "..\branding\installer_header.bmp"
!define MUI_WELCOMEFINISHPAGE_BITMAP "..\branding\installer_welcome.bmp"
!define MUI_UNWELCOMEFINISHPAGE_BITMAP "..\branding\installer_welcome.bmp"
!define MUI_BGCOLOR "FBF5EC"
!define MUI_TEXTCOLOR "2E2A33"
!define MUI_INSTFILESPAGE_COLORS "2E2A33 FBF5EC"
!define MUI_INSTFILESPAGE_PROGRESSBAR "smooth"
!define MUI_ABORTWARNING
!define MUI_ABORTWARNING_TEXT "Are you sure you want to stop installing Snugfit?"

; ---------------------------------------------------------------- pages
!define MUI_WELCOMEPAGE_TITLE "Welcome to Snugfit!"
!define MUI_WELCOMEPAGE_TEXT "Snugfit is a cozy jigsaw puzzle game for Windows.$\r$\n$\r$\n\
Version ${VERSION} comes with 10 beautiful puzzles, including the October Collection, and four difficulty levels from 100 to 1,000 pieces.$\r$\n$\r$\n\
More puzzles are coming soon: five new ones every month, delivered by the built-in updater.$\r$\n$\r$\n\
Click Next to continue."
!insertmacro MUI_PAGE_WELCOME

!define MUI_DIRECTORYPAGE_TEXT_TOP "Snugfit will be installed for your Windows account in the folder below. No administrator rights are needed."
!insertmacro MUI_PAGE_DIRECTORY

!insertmacro MUI_PAGE_INSTFILES

!define MUI_FINISHPAGE_TITLE "Snugfit is ready!"
!define MUI_FINISHPAGE_TEXT "Your puzzles are waiting. Your saved games and settings are kept safe between updates.$\r$\n$\r$\nHappy puzzling!"
!define MUI_FINISHPAGE_RUN "$INSTDIR\Snugfit.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Start Snugfit now"
!define MUI_FINISHPAGE_SHOWREADME ""
!define MUI_FINISHPAGE_SHOWREADME_TEXT "Create a desktop shortcut"
!define MUI_FINISHPAGE_SHOWREADME_FUNCTION CreateDesktopShortcut
!define MUI_FINISHPAGE_LINK "Snugfit on GitHub"
!define MUI_FINISHPAGE_LINK_LOCATION "${WEBSITE}"
!insertmacro MUI_PAGE_FINISH

!define MUI_WELCOMEPAGE_TITLE "Uninstall Snugfit"
!define MUI_WELCOMEPAGE_TEXT "This will remove Snugfit from your computer.$\r$\n$\r$\nClick Next to continue."
!insertmacro MUI_UNPAGE_WELCOME
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_UNPAGE_FINISH

!insertmacro MUI_LANGUAGE "English"

; ---------------------------------------------------------------- functions
Function .onInit
  ; themed splash screen
  InitPluginsDir
  File "/oname=$PLUGINSDIR\splash.bmp" "..\branding\installer_splash.bmp"
  IfSilent +2
  advsplash::show 1400 300 250 -1 "$PLUGINSDIR\splash"
  Pop $0
FunctionEnd

Function CreateDesktopShortcut
  CreateShortcut "$DESKTOP\Snugfit.lnk" "$INSTDIR\Snugfit.exe" "" "$INSTDIR\Snugfit.exe" 0
FunctionEnd

; ---------------------------------------------------------------- install
Section "Snugfit" SecMain
  SectionIn RO
  ; close a running copy (e.g. when updating from inside the app)
  nsExec::Exec 'taskkill /IM Snugfit.exe /F'
  Pop $0
  Sleep 400

  SetOutPath "$INSTDIR"
  ; replace bundled assets cleanly (user saves live in %APPDATA%\Snugfit, not here)
  RMDir /r "$INSTDIR\assets"
  File "${SRCDIR}\Snugfit.exe"
  File /nonfatal "${SRCDIR}\LICENSE.txt"
  File /nonfatal "${SRCDIR}\CHANGELOG.md"
  File /nonfatal "${SRCDIR}\THIRD-PARTY-NOTICES.txt"
  SetOutPath "$INSTDIR\assets"
  File /r "${SRCDIR}\assets\*.*"
  SetOutPath "$INSTDIR"

  WriteUninstaller "$INSTDIR\Uninstall Snugfit.exe"
  CreateShortcut "$SMPROGRAMS\Snugfit.lnk" "$INSTDIR\Snugfit.exe" "" "$INSTDIR\Snugfit.exe" 0

  WriteRegStr HKCU "Software\Snugfit" "InstallDir" "$INSTDIR"
  WriteRegStr HKCU "${UNINSTKEY}" "DisplayName" "Snugfit"
  WriteRegStr HKCU "${UNINSTKEY}" "DisplayVersion" "${VERSION}"
  WriteRegStr HKCU "${UNINSTKEY}" "Publisher" "${PUBLISHER}"
  WriteRegStr HKCU "${UNINSTKEY}" "DisplayIcon" "$INSTDIR\Snugfit.exe"
  WriteRegStr HKCU "${UNINSTKEY}" "InstallLocation" "$INSTDIR"
  WriteRegStr HKCU "${UNINSTKEY}" "URLInfoAbout" "${WEBSITE}"
  WriteRegStr HKCU "${UNINSTKEY}" "UninstallString" '"$INSTDIR\Uninstall Snugfit.exe"'
  WriteRegStr HKCU "${UNINSTKEY}" "QuietUninstallString" '"$INSTDIR\Uninstall Snugfit.exe" /S'
  WriteRegDWORD HKCU "${UNINSTKEY}" "NoModify" 1
  WriteRegDWORD HKCU "${UNINSTKEY}" "NoRepair" 1
  ${GetSize} "$INSTDIR" "/S=0K" $0 $1 $2
  IntFmt $0 "0x%08X" $0
  WriteRegDWORD HKCU "${UNINSTKEY}" "EstimatedSize" "$0"
SectionEnd

; ---------------------------------------------------------------- uninstall
Section "Uninstall"
  nsExec::Exec 'taskkill /IM Snugfit.exe /F'
  Pop $0
  Delete "$SMPROGRAMS\Snugfit.lnk"
  Delete "$DESKTOP\Snugfit.lnk"
  Delete "$INSTDIR\Snugfit.exe"
  Delete "$INSTDIR\LICENSE.txt"
  Delete "$INSTDIR\CHANGELOG.md"
  Delete "$INSTDIR\THIRD-PARTY-NOTICES.txt"
  RMDir /r "$INSTDIR\assets"
  Delete "$INSTDIR\Uninstall Snugfit.exe"
  RMDir "$INSTDIR"
  DeleteRegKey HKCU "${UNINSTKEY}"
  DeleteRegKey HKCU "Software\Snugfit"

  IfSilent done
  MessageBox MB_YESNO|MB_ICONQUESTION|MB_DEFBUTTON2 "Also delete your saved puzzles and settings?$\r$\n$\r$\nChoose No to keep them in case you reinstall Snugfit later." IDNO done
  RMDir /r "$APPDATA\Snugfit"
  done:
SectionEnd
