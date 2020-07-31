; S3 Explorer
; Copyright (C) 2020  indece UG (haftungsbeschränkt)
;
; This program is free software: you can redistribute it and/or modify
; it under the terms of the GNU General Public License as published by
; the Free Software Foundation, either version 3 of the License or any
; later version.
;
; This program is distributed in the hope that it will be useful,
; but WITHOUT ANY WARRANTY; without even the implied warranty of
; MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
; GNU General Public License for more details.
;
; You should have received a copy of the GNU General Public License
; along with this program. If not, see <https://www.gnu.org/licenses/>.

; Script generated by the Inno Setup Script Wizard.
; SEE THE DOCUMENTATION FOR DETAILS ON CREATING INNO SETUP SCRIPT FILES!
#define MyAppName "S3 Explorer"
#define MyAppVersion GetEnv('BUILD_VERSION')
#define MyAppPublisher "indece UG (haftungsbeschränkt)"
#define MyAppURL "https://www.indece.com/software/s3-explorer"
#define MyAppExeName "s3-explorer.exe"

[Setup]
; NOTE: The value of AppId uniquely identifies this application. Do not use the same AppId value in installers for other applications.
; (To generate a new GUID, click Tools | Generate GUID inside the IDE.)
AppId={{61580A48-CFEE-4494-A06B-ADC65D6D5C43}
AppName={#MyAppName}
AppVersion={#MyAppVersion}
;AppVerName={#MyAppName} {#MyAppVersion}
AppPublisher={#MyAppPublisher}
AppPublisherURL={#MyAppURL}
AppSupportURL={#MyAppURL}
AppUpdatesURL={#MyAppURL}
DefaultDirName={autopf}\S3 Explorer
DisableProgramGroupPage=yes
LicenseFile=assets/bin/LICENSE.txt
OutputDir=dist/
OutputBaseFilename=s3-explorer-setup_win64-unsigned
Compression=lzma
SolidCompression=yes
WizardStyle=modern

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "german"; MessagesFile: "compiler:Languages\German.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "assets/bin/s3-explorer_win64.exe"; DestDir: "{app}"; DestName: "s3-explorer.exe"; Flags: ignoreversion
Source: "assets/bin/LICENSE.txt"; DestDir: "{app}"; Flags: ignoreversion
Source: "assets/bin/licenses.zip"; DestDir: "{app}"; Flags: ignoreversion
Source: "assets/webview.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "assets/WebView2Loader.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "assets/s3-explorer.png"; DestDir: "{app}"; Flags: ignoreversion
; NOTE: Don't use "Flags: ignoreversion" on any shared system files

[Icons]
Name: "{autoprograms}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"
Name: "{autodesktop}\{#MyAppName}"; Filename: "{app}\{#MyAppExeName}"; Tasks: desktopicon

[Run]
Filename: "{app}\{#MyAppExeName}"; Description: "{cm:LaunchProgram,{#StringChange(MyAppName, '&', '&&')}}"; Flags: nowait postinstall skipifsilent

[Code]

procedure CurStepChanged(CurStep: TSetupStep);
var
  Code: Integer;
begin
  if CurStep = ssInstall then
  begin
    Exec('CheckNetIsolation.exe', 'LoopbackExempt -a -n="Microsoft.Win32WebViewHost_cw5n1h2txyewy"', '', SW_HIDE, ewWaitUntilTerminated, Code);
  end;
end;