-- Launch adl exactly as if it were typed in the active Terminal tab.
on run
	tell application "Terminal"
		if (count of windows) is 0 then
			do script "adl"
		else
			set targetTab to selected tab of front window

			-- Do not interrupt a command that is already running. In that case,
			-- open a new window and copy the active tab's appearance.
			if busy of targetTab then
				set activeSettings to current settings of targetTab
				set targetTab to do script ""
				set current settings of targetTab to activeSettings
			end if

			do script "adl" in targetTab
		end if

		activate
	end tell
end run
