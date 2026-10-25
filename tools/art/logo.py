"""Snugfit logo, icon and installer graphics."""
import math
from common import *

ORANGE = '#F07A2A'
TEAL = '#2C9A92'
PLUM = '#5B3F8C'
INK = '#2E2A33'
CREAM = '#FBF5EC'

TAB = [(0.35, 0), 'C', (0.37, 0), (0.39, 0.02), (0.38, 0.06), 'C', (0.35, 0.13), (0.38, 0.23), (0.5, 0.23),
       'C', (0.62, 0.23), (0.65, 0.13), (0.62, 0.06), 'C', (0.61, 0.02), (0.63, 0), (0.65, 0), 'L', (1, 0)]


def edge(p0, p1, kind):
    dx, dy = p1[0]-p0[0], p1[1]-p0[1]
    nx, ny = dy, -dx  # outward for clockwise traversal (y down)
    def m(t, n):
        n *= kind
        return (p0[0] + dx*t + nx*n, p0[1] + dy*t + ny*n)
    if kind == 0:
        return ' L%.2f %.2f' % p1
    out = ' L%.2f %.2f' % m(*TAB[0])
    i = 1
    while i < len(TAB):
        op = TAB[i]
        if op == 'C':
            out += ' C' + ' '.join('%.2f %.2f' % m(*TAB[i+k]) for k in (1, 2, 3)); i += 4
        else:
            out += ' L%.2f %.2f' % m(*TAB[i+1]); i += 2
    return out


def piece(x, y, s, top, right, bottom, left):
    a, b, c, d = (x, y), (x+s, y), (x+s, y+s), (x, y+s)
    return 'M%.2f %.2f' % a + edge(a, b, top) + edge(b, c, right) + edge(c, d, bottom) + edge(d, a, left) + 'Z'


def pieces_group(cx, cy, s, rot=-10, stroke=None, sw=0):
    pa = piece(-s, -s/2, s, 1, 1, -1, 0)
    pb = piece(0, -s/2, s, -1, 0, 1, -1)
    st = ('stroke="%s" stroke-width="%.1f" stroke-linejoin="round"' % (stroke, sw)) if stroke else ''
    hl = ('<path d="%s" fill="none" stroke="#fff" stroke-opacity="0.35" stroke-width="%.1f" transform="translate(%.1f %.1f)" clip-path="url(#ca)"/>')
    return ('<g transform="translate(%.1f %.1f) rotate(%.1f)">'
            '<path d="%s" fill="#000" opacity="0.18" transform="translate(%.1f %.1f)"/>'
            '<path d="%s" fill="#000" opacity="0.18" transform="translate(%.1f %.1f)"/>'
            '<path d="%s" fill="url(#gOrange)" %s/>'
            '<path d="%s" fill="url(#gTeal)" %s/>'
            '<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#fff" opacity="0.9"/>'
            '</g>') % (cx, cy, rot, pa, s*0.04, s*0.06, pb, s*0.04, s*0.06, pa, st, pb, st, s*0.62, -s*0.12, s*0.07)


GRADS = '''<linearGradient id="gOrange" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#FF9A4A"/><stop offset="1" stop-color="#E2601B"/></linearGradient>
<linearGradient id="gTeal" x1="0" y1="0" x2="1" y2="1"><stop offset="0" stop-color="#3DBDB3"/><stop offset="1" stop-color="#1F7C76"/></linearGradient>
<linearGradient id="gPlum" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#6E4FA6"/><stop offset="1" stop-color="#3C2862"/></linearGradient>'''


def icon_svg(size=1024, bg=True):
    s = size
    body = ''
    if bg:
        body += '<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="%.1f" fill="%s"/>' % (s*0.04, s*0.04, s*0.92, s*0.92, s*0.22, CREAM)
        body += '<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="%.1f" fill="none" stroke="#EADFCC" stroke-width="%.1f"/>' % (s*0.04, s*0.04, s*0.92, s*0.92, s*0.22, s*0.012)
    body += pieces_group(s*0.5, s*0.53, s*0.33, -10, '#ffffff', s*0.018)
    return '<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d"><defs>%s</defs>%s</svg>' % (s, s, s, s, GRADS, body)


def wordmark_svg(w=1400, h=420, ink=INK, sub=True):
    body = pieces_group(210, 210, 120, -10, '#ffffff', 7)
    body += '<text x="420" y="270" font-family="Poppins" font-weight="bold" font-size="210" fill="%s" letter-spacing="-4">Snug<tspan fill="%s">fit</tspan></text>' % (ink, ORANGE)
    if sub:
        body += '<text x="428" y="345" font-family="Poppins" font-size="52" fill="%s" opacity="0.65" letter-spacing="10">JIGSAW PUZZLES</text>' % ink
    return '<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d"><defs>%s</defs>%s</svg>' % (w, h, w, h, GRADS, body)


def welcome_svg(w=164*3, h=314*3):
    body = '<rect width="%d" height="%d" fill="url(#gPlum)"/>' % (w, h)
    import random
    rng = random.Random(3)
    for i in range(14):
        x, y = rng.uniform(0, w), rng.uniform(0, h)
        s = rng.uniform(40, 90)
        col = rng.choice(['#ffffff', ORANGE, TEAL, '#F6C453'])
        body += '<path d="%s" fill="%s" opacity="%.2f" transform="translate(%.0f %.0f) rotate(%.0f)"/>' % (piece(-s/2, -s/2, s, 1, -1, 1, -1), col, rng.uniform(0.08, 0.2), x, y, rng.uniform(0, 360))
    body += '<circle cx="%d" cy="%d" r="%d" fill="%s"/>' % (w/2, h*0.36, w*0.36, CREAM)
    body += pieces_group(w/2, h*0.37, w*0.16, -10, '#ffffff', 7)
    body += '<text x="%d" y="%d" font-family="Poppins" font-weight="bold" font-size="96" fill="#ffffff" text-anchor="middle">Snug<tspan fill="#FFB074">fit</tspan></text>' % (w/2, h*0.68)
    body += '<text x="%d" y="%d" font-family="Poppins" font-size="30" fill="#ffffff" opacity="0.75" text-anchor="middle" letter-spacing="6">JIGSAW PUZZLES</text>' % (w/2, h*0.68 + 56)
    body += '<text x="%d" y="%d" font-family="Poppins" font-size="26" fill="#ffffff" opacity="0.6" text-anchor="middle">New puzzles every month</text>' % (w/2, h*0.93)
    return '<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d"><defs>%s</defs>%s</svg>' % (w, h, w, h, GRADS, body)


def header_svg(w=150*3, h=57*3):
    body = '<rect width="%d" height="%d" fill="#FFFFFF"/>' % (w, h)
    body += pieces_group(80, h/2 + 6, 46, -10, '#ffffff', 4)
    body += '<text x="160" y="%d" font-family="Poppins" font-weight="bold" font-size="66" fill="%s">Snug<tspan fill="%s">fit</tspan></text>' % (h/2 + 22, INK, ORANGE)
    return '<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d"><defs>%s</defs>%s</svg>' % (w, h, w, h, GRADS, body)


def render_transparent(svg, out, w, h):
    from playwright.sync_api import sync_playwright
    with sync_playwright() as p:
        b = p.chromium.launch()
        pg = b.new_page(viewport={'width': w, 'height': h})
        pg.set_content('<html><body style="margin:0;background:transparent">%s</body></html>' % svg)
        pg.wait_for_timeout(300)
        pg.screenshot(path=out, clip={'x': 0, 'y': 0, 'width': w, 'height': h}, omit_background=True)
        b.close()


if __name__ == '__main__':
    import sys, os
    from PIL import Image
    out = sys.argv[1]
    os.makedirs(out, exist_ok=True)
    render_transparent(icon_svg(1024), out + '/icon_1024.png', 1024, 1024)
    render_transparent(wordmark_svg(), out + '/logo_light.png', 1400, 420)
    render_transparent(wordmark_svg(ink='#F4EEE6'), out + '/logo_dark.png', 1400, 420)
    render_transparent(welcome_svg(), out + '/_welcome.png', 164*3, 314*3)
    render_transparent(header_svg(), out + '/_header.png', 150*3, 57*3)
    ic = Image.open(out + '/icon_1024.png').convert('RGBA')
    ic.save(out + '/snugfit.ico', sizes=[(16, 16), (24, 24), (32, 32), (48, 48), (64, 64), (128, 128), (256, 256)])
    ic.resize((256, 256), Image.LANCZOS).save(out + '/icon_256.png')
    ic.resize((64, 64), Image.LANCZOS).save(out + '/icon_64.png')
    ic.resize((32, 32), Image.LANCZOS).save(out + '/icon_32.png')
    ic.resize((16, 16), Image.LANCZOS).save(out + '/icon_16.png')
    Image.open(out + '/_welcome.png').convert('RGB').resize((164, 314), Image.LANCZOS).save(out + '/installer_welcome.bmp')
    Image.open(out + '/_header.png').convert('RGB').resize((150, 57), Image.LANCZOS).save(out + '/installer_header.bmp')
    print('ok')
