# Microsoft Defender False Positive Submission Guide

If `konfetti-win.exe` is flagged (e.g., `Win32/Sabsik.FL.A!ml`), follow this process to reduce future detections.

## 1. Verify Clean Build
- Rebuild on a clean machine/container
```bash
make win
```
- Do NOT pack or compress the binary (avoid UPX, etc.)
- Run an external scan (optional):
  - https://www.virustotal.com

## 2. Produce Hashes
Generate common hashes for identification:
```powershell
Get-FileHash .\dist\konfetti-win.exe -Algorithm SHA256
Get-FileHash .\dist\konfetti-win.exe -Algorithm SHA1
```
Or on mac/Linux:
```bash
shasum -a 256 dist/konfetti-win.exe
shasum -a 1 dist/konfetti-win.exe
```

## 3. Submit to Microsoft
Go to: https://www.microsoft.com/en-us/wdsi/filesubmission

Fill out:
- Submission type: Software
- Classification: False Positive
- File: `konfetti-win.exe`
- Description: "Open source config inspection tool built in Go. No networking, no persistence, no self-modifying behavior. Source: https://github.com/bird2920/Konfetti"
- Publisher: bird2920
- Detection name: (e.g.) `Win32/Sabsik.FL.A!ml`

## 4. Optional: Sign the Binary
Acquire a code-signing certificate (DigiCert, Sectigo, etc.). Sign:
```powershell
signtool sign /fd SHA256 /tr http://timestamp.digicert.com /td SHA256 /a dist/konfetti-win.exe
```
Verify signature:
```powershell
signtool verify /pa dist/konfetti-win.exe
```

## 5. Re-Scan
After Defender definitions update (can take hours–days), re-download and test on a fresh VM.

## 6. Communicate to Users
Add a README note:
> Windows Defender may briefly flag fresh unsigned builds. If you see `Win32/Sabsik.FL.A!ml`, it's a generic ML heuristic. Submit the binary as a false positive or use the signed release.

## 7. Deterministic Builds (Advanced)
Use repeatable flags to keep binaries stable between releases:
```bash
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -X main.Version=2.0.1" -o dist/konfetti-win.exe .
```
This removes symbol/debug info (-s -w) and strips local paths (trimpath) to reduce variance.

## 8. Track Changes
Maintain a release log of SHA256 hashes for each published version:
```
Version 2.0.1
SHA256: <paste-hash>
Built: 2025-10-28
Go: go1.xx
```

## 9. If Still Flagged
- Open a support ticket with Microsoft referencing previous submission ID.
- Provide build command, full SHA256, and VirusTotal results.

---
Prepared for Konfetti by bird2920.
