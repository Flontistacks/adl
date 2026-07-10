-- Launch adl in Terminal
--
-- Script Editor:
-- 1. Open in Script Editor
-- 2. File > Export > Application
-- 3. Save as "adl" and drag to Dock

on run
	set adlPath to locateAdl()

	if adlPath is "" then
		display dialog "adl was not found on this Mac." & return & return & "Install it with:" & return & "brew install Flontistacks/tap/adl" buttons {"OK"} default button 1 with icon caution with title "adl"
		return
	end if

	set shellCmd to "export PATH=/opt/homebrew/bin:/usr/local/bin:/usr/bin:/bin:$PATH; " & quoted form of adlPath & "; exit"

	tell application "Terminal"
		activate
		do script shellCmd
	end tell
end run

on locateAdl()
	if pathExists("/opt/homebrew/bin/adl") then return "/opt/homebrew/bin/adl"
	if pathExists("/usr/local/bin/adl") then return "/usr/local/bin/adl"
	return ""
end locateAdl

on pathExists(posixPath)
	try
		do shell script "test -f " & quoted form of posixPath
		return true
	on error
		return false
	end try
end pathExists
