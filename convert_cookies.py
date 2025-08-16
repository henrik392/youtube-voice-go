def cookies_to_string(file_path: str) -> str:
    with open(file_path, "r", encoding="utf-8") as f:
        lines = f.readlines()

    processed_lines = []
    for line in lines:
        line = line.rstrip('\n\r')  # Remove line endings
        
        # If it's a comment line, keep it as-is
        if line.startswith('#') or line.strip() == '':
            processed_lines.append(line)
        else:
            # For cookie data lines, replace multiple spaces/tabs with single tab
            import re
            # Split on whitespace and rejoin with tabs
            parts = re.split(r'\s+', line.strip())
            if len(parts) >= 7:  # Valid cookie line should have at least 7 fields
                processed_lines.append('\t'.join(parts))
            else:
                processed_lines.append(line)  # Keep invalid lines as-is

    # Join with literal \n and wrap in quotes
    content = '\\n'.join(processed_lines)
    return f'"{content}"'


if __name__ == "__main__":
    file_path = "cookies.txt"  # change if needed
    result = cookies_to_string(file_path)
    print(result)
