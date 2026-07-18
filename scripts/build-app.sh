#!/bin/bash
# Builds adl.app — drag to Applications, set your own icon via Get Info.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
APP="$SCRIPT_DIR/adl.app"

rm -rf "$APP"
mkdir -p "$APP/Contents/MacOS" "$APP/Contents/Resources"

cp "$SCRIPT_DIR/Launch adl.applescript" "$APP/Contents/Resources/launch.applescript"

cat > "$APP/Contents/MacOS/adl" << 'EOF'
#!/bin/bash
DIR="$(cd "$(dirname "$0")/../Resources" && pwd)"
exec /usr/bin/osascript "$DIR/launch.applescript"
EOF
chmod +x "$APP/Contents/MacOS/adl"

cat > "$APP/Contents/Info.plist" << 'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleExecutable</key>
	<string>adl</string>
	<key>CFBundleIdentifier</key>
	<string>com.flontistacks.adl</string>
	<key>CFBundleName</key>
	<string>adl</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>1.0</string>
	<key>CFBundleVersion</key>
	<string>1</string>
	<key>LSMinimumSystemVersion</key>
	<string>12.0</string>
	<key>NSHighResolutionCapable</key>
	<true/>
</dict>
</plist>
EOF

echo "Built: $APP"
echo "Copy to Applications (replace old app), keep your icon, drag to Dock."
