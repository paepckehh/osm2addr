#!/usr/bin/env python3
"""Render captured osm2addr output as a terminal-style SVG screenshot.

Usage: render_terminal_svg.py INPUT OUTPUT TITLE COMMAND [light|dark]
"""
import html
import re
import sys

THEMES = {
    "dark": {
        "bg0": "#0d1728",
        "bg1": "#091120",
        "bar0": "#122037",
        "bar1": "#0f1b2e",
        "stroke": "#22344f",
        "rule": "#1e2e47",
        "title": "#7488a3",
        "text": "#c9d6e8",
        "mute": "#7d8fa8",
        "sep": "#3b506e",
        "accent": "#5eead4",
        "blue": "#a8cbff",
        "amber": "#f8d78a",
        "purple": "#c4b5fd",
        "prompt": "#34d399",
        "shadow": "#000000",
        "shadow_opacity": "0.45",
    },
    "light": {
        "bg0": "#ffffff",
        "bg1": "#f4f8fd",
        "bar0": "#eef4fb",
        "bar1": "#e5eefb",
        "stroke": "#d7e3f2",
        "rule": "#dbe7f5",
        "title": "#64748b",
        "text": "#1e293b",
        "mute": "#64748b",
        "sep": "#aabdd4",
        "accent": "#0f766e",
        "blue": "#1d4ed8",
        "amber": "#b45309",
        "purple": "#7c3aed",
        "prompt": "#059669",
        "shadow": "#7f9cc0",
        "shadow_opacity": "0.28",
    },
}

JSON_PAIR = re.compile(r'^(\s*)("(?:[^"\\]|\\.)*")(\s*:\s*)("(?:[^"\\]|\\.)*")(,?)$')
JSON_EDGE = re.compile(r'^(\s*)([{}])(,?)$')
HEX_ID = re.compile(r'^"[0-9a-f]{24}"$')


def esc(text):
    return html.escape(text, quote=False)


def value_color(key, theme):
    if "Err:" in key:
        return theme["amber"]
    if "Corrected:" in key:
        return theme["purple"]
    if "Writer:" in key or "File" in key or "URL" in key:
        return theme["blue"]
    if "Collect:" in key or "Time:" in key or "Complete:" in key:
        return theme["accent"]
    return theme["text"]


def report_line(line, theme):
    if line.startswith("---"):
        return None
    key, sep, value = line.partition("#")
    if not sep:
        return f'<tspan fill="{theme["text"]}">{esc(line)}</tspan>'
    return (
        f'<tspan fill="{theme["mute"]}">{esc(key)}</tspan>'
        f'<tspan fill="{theme["sep"]}">#</tspan>'
        f'<tspan fill="{value_color(key, theme)}">{esc(value)}</tspan>'
    )


def json_line(line, theme):
    pair = JSON_PAIR.match(line)
    if pair:
        indent, key, colon, value, comma = pair.groups()
        color = theme["amber"] if HEX_ID.match(value) else theme["accent"]
        return (
            f'{esc(indent)}<tspan fill="{theme["blue"]}">{esc(key)}</tspan>'
            f'<tspan fill="{theme["sep"]}">{esc(colon)}</tspan>'
            f'<tspan fill="{color}">{esc(value)}</tspan>'
            f'<tspan fill="{theme["mute"]}">{esc(comma)}</tspan>'
        )
    edge = JSON_EDGE.match(line)
    if edge:
        indent, brace, comma = edge.groups()
        return (
            f'{esc(indent)}<tspan fill="{theme["mute"]}">{esc(brace + comma)}</tspan>'
        )
    return f'<tspan fill="{theme["text"]}">{esc(line)}</tspan>'


def render(content, outfile, title, command, theme_name):
    theme = THEMES[theme_name]
    lines = [line.expandtabs(2) for line in content.splitlines()]
    while lines and not lines[-1].strip():
        lines.pop()
    height = 110 + len(lines) * 21 + 22
    is_json = lines and lines[0].lstrip().startswith("{")
    body = []
    y = 110
    for line in lines:
        if not is_json and line.startswith("---"):
            body.append(
                f'    <line x1="32" y1="{y - 5}" x2="968" y2="{y - 5}" '
                f'stroke="{theme["rule"]}"/>'
            )
        else:
            markup = json_line(line, theme) if is_json else report_line(line, theme)
            body.append(
                f'    <text x="32" y="{y}" xml:space="preserve">{markup}</text>'
            )
        y += 21
    svg = f'''<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 1000 {height}" width="1000" height="{height}" role="img" aria-label="{esc(title)}">
  <title>{esc(title)}</title>
  <defs>
    <linearGradient id="term" x1="0" y1="0" x2="0" y2="1">
      <stop offset="0" stop-color="{theme["bg0"]}"/>
      <stop offset="1" stop-color="{theme["bg1"]}"/>
    </linearGradient>
    <linearGradient id="bar" x1="0" y1="0" x2="1" y2="0">
      <stop offset="0" stop-color="{theme["bar0"]}"/>
      <stop offset="1" stop-color="{theme["bar1"]}"/>
    </linearGradient>
    <filter id="soft" x="-10%" y="-10%" width="120%" height="130%">
      <feDropShadow dx="0" dy="12" stdDeviation="16" flood-color="{theme["shadow"]}" flood-opacity="{theme["shadow_opacity"]}"/>
    </filter>
    <style>
      .mono {{ font-family: ui-monospace, SFMono-Regular, "SF Mono", Menlo, Consolas, "Liberation Mono", monospace; }}
    </style>
  </defs>
  <g filter="url(#soft)">
    <rect width="1000" height="{height}" rx="16" fill="url(#term)" stroke="{theme["stroke"]}"/>
  </g>
  <path d="M0 16 A16 16 0 0 1 16 0 H984 A16 16 0 0 1 1000 16 V44 H0 Z" fill="url(#bar)"/>
  <line x1="0" y1="44.5" x2="1000" y2="44.5" stroke="{theme["rule"]}"/>
  <circle cx="24" cy="22" r="6" fill="#f87171"/>
  <circle cx="44" cy="22" r="6" fill="#fbbf24"/>
  <circle cx="64" cy="22" r="6" fill="#34d399"/>
  <text x="500" y="27" class="mono" font-size="12" text-anchor="middle" fill="{theme["title"]}">{esc(title)}</text>
  <text x="32" y="72" class="mono" font-size="14" xml:space="preserve"><tspan fill="{theme["prompt"]}">$ </tspan><tspan fill="{theme["text"]}">{esc(command)}</tspan></text>
  <line x1="0" y1="86" x2="1000" y2="86" stroke="{theme["rule"]}"/>
  <g class="mono" font-size="14">
{chr(10).join(body)}
  </g>
</svg>
'''
    with open(outfile, "w", encoding="utf-8") as handle:
        handle.write(svg)
    print(f"wrote {outfile} ({height}px)")


def main():
    if len(sys.argv) not in (5, 6):
        sys.exit("usage: render_terminal_svg.py INPUT OUTPUT TITLE COMMAND [light|dark]")
    infile, outfile, title, command = sys.argv[1:5]
    theme_name = sys.argv[5] if len(sys.argv) == 6 else "dark"
    if theme_name not in THEMES:
        sys.exit(f"unknown theme {theme_name!r}, expected light or dark")
    with open(infile, encoding="utf-8") as handle:
        render(handle.read(), outfile, title, command, theme_name)


if __name__ == "__main__":
    main()
