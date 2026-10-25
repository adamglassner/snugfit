"""Everyday puzzle scenes: Boston Terriers, flowers, ocean, forest, stained glass."""
import math, random
from common import *


# ---------------------------------------------------------------- Boston Terriers
def boston_head(cx, cy, s, rng, coat='#1d1c22', acc=None, tilt=0, bg='#fff'):
    """Front-facing cartoon Boston Terrier, head width ~= 2*s"""
    o = ['<g transform="translate(%.0f %.0f) rotate(%.1f) scale(%.2f %.2f)">' % (cx, cy, tilt, s*1.12, s)]
    white = '#f7f4ef'
    pink = '#e9a0a6'
    hi = shade(coat, 0.18)
    # body / chest
    o.append('<path d="M-0.95 1.9 C-0.95 1.1 -0.6 0.75 0 0.75 C0.6 0.75 0.95 1.1 0.95 1.9Z" fill="%s"/>' % coat)
    o.append('<path d="M-0.42 1.9 C-0.4 1.3 -0.25 0.9 0 0.85 C0.25 0.9 0.4 1.3 0.42 1.9Z" fill="%s"/>' % white)
    # ears
    for sgn in (-1, 1):
        o.append('<path d="M%.2f -0.35 C%.2f -0.9 %.2f -1.35 %.2f -1.45 C%.2f -1.25 %.2f -0.85 %.2f -0.45Z" fill="%s"/>' % (
            sgn*0.18, sgn*0.35, sgn*0.6, sgn*0.82, sgn*0.9, sgn*0.85, sgn*0.78, coat))
        o.append('<path d="M%.2f -0.48 C%.2f -0.9 %.2f -1.18 %.2f -1.26 C%.2f -1.1 %.2f -0.8 %.2f -0.55Z" fill="%s"/>' % (
            sgn*0.32, sgn*0.42, sgn*0.62, sgn*0.77, sgn*0.8, sgn*0.76, sgn*0.7, pink))
    # head
    o.append('<path d="M-0.82 -0.1 C-0.86 -0.62 -0.45 -0.85 0 -0.85 C0.45 -0.85 0.86 -0.62 0.82 -0.1 C0.8 0.35 0.55 0.72 0 0.74 C-0.55 0.72 -0.8 0.35 -0.82 -0.1Z" fill="%s"/>' % coat)
    o.append('<ellipse cx="-0.35" cy="-0.55" rx="0.22" ry="0.1" fill="%s" opacity="0.5"/>' % hi)
    # blaze + muzzle
    o.append('<path d="M-0.07 -0.84 C-0.05 -0.6 -0.1 -0.35 -0.18 -0.15 C-0.42 0.0 -0.48 0.3 -0.38 0.5 C-0.25 0.72 0.25 0.72 0.38 0.5 C0.48 0.3 0.42 0.0 0.18 -0.15 C0.1 -0.35 0.05 -0.6 0.07 -0.84Z" fill="%s"/>' % white)
    # eyes
    for sgn in (-1, 1):
        o.append('<circle cx="%.2f" cy="-0.18" r="0.17" fill="#fff" opacity="0.15"/>' % (sgn*0.4))
        o.append('<circle cx="%.2f" cy="-0.18" r="0.15" fill="#3b2414"/><circle cx="%.2f" cy="-0.18" r="0.1" fill="#120a06"/>' % (sgn*0.4, sgn*0.4))
        o.append('<circle cx="%.2f" cy="-0.24" r="0.045" fill="#fff"/><circle cx="%.2f" cy="-0.12" r="0.02" fill="#fff"/>' % (sgn*0.4 - 0.04, sgn*0.4 + 0.05))
    # nose
    o.append('<path d="M-0.13 0.1 C-0.14 0.02 0.14 0.02 0.13 0.1 C0.12 0.2 0.04 0.24 0 0.24 C-0.04 0.24 -0.12 0.2 -0.13 0.1Z" fill="#141216"/>')
    o.append('<ellipse cx="-0.04" cy="0.08" rx="0.04" ry="0.02" fill="#fff" opacity="0.6"/>')
    # mouth
    tongue = rng.random() < 0.55
    if tongue:
        o.append('<path d="M-0.08 0.36 C-0.1 0.52 0.1 0.52 0.08 0.36Z" fill="#e66b7d"/>')
    o.append('<path d="M0 0.24 L0 0.33 M-0.2 0.3 C-0.12 0.4 -0.04 0.38 0 0.33 C0.04 0.38 0.12 0.4 0.2 0.3" stroke="#141216" stroke-width="0.03" fill="none" stroke-linecap="round"/>')
    for sgn in (-1, 1):
        for k in range(3):
            o.append('<circle cx="%.2f" cy="%.2f" r="0.012" fill="#555"/>' % (sgn*(0.14 + k*0.06), 0.24 + (k % 2)*0.05))
    # accessory
    if acc == 'bandana':
        c = rng.choice(['#d63031', '#0984e3', '#e17055', '#00b894', '#fdcb6e', '#6c5ce7'])
        o.append('<path d="M-0.62 0.72 C-0.3 0.9 0.3 0.9 0.62 0.72 L0 1.35Z" fill="%s"/>' % c)
        for k in range(6):
            o.append('<circle cx="%.2f" cy="%.2f" r="0.04" fill="#fff" opacity="0.85"/>' % (-0.35 + k*0.14, 0.9 + (k % 2)*0.12))
    elif acc == 'bowtie':
        c = rng.choice(['#d63031', '#2d3436', '#0984e3', '#e84393'])
        o.append('<path d="M0 0.88 L-0.3 0.72 L-0.3 1.04Z M0 0.88 L0.3 0.72 L0.3 1.04Z" fill="%s"/><circle cx="0" cy="0.88" r="0.08" fill="%s"/>' % (c, shade(c, -0.2)))
    elif acc == 'collar':
        c = rng.choice(['#e17055', '#00cec9', '#a29bfe', '#fab1a0'])
        o.append('<path d="M-0.6 0.72 C-0.3 0.86 0.3 0.86 0.6 0.72 L0.6 0.84 C0.3 0.98 -0.3 0.98 -0.6 0.84Z" fill="%s"/><circle cx="0" cy="1.0" r="0.1" fill="#f1c40f" stroke="#b7950b" stroke-width="0.02"/>' % c)
    elif acc == 'party':
        c = rng.choice(['#e84393', '#00b894', '#fdcb6e', '#74b9ff'])
        o.append('<g transform="rotate(14)"><path d="M-0.25 -0.78 L0.25 -0.78 L0 -1.5Z" fill="%s"/><circle cx="0" cy="-1.52" r="0.07" fill="#fff"/><path d="M-0.17 -1.0 L0.17 -1.0 M-0.1 -1.2 L0.1 -1.2" stroke="#fff" stroke-width="0.05"/></g>' % c)
    elif acc == 'flowers':
        for k in range(5):
            fx = -0.5 + k*0.25
            fy = -0.8 - math.sin(k/4*math.pi)*0.12
            fc = rng.choice(['#ff7675', '#fdcb6e', '#fd79a8', '#a29bfe', '#fff'])
            for a in range(5):
                o.append('<circle cx="%.2f" cy="%.2f" r="0.06" fill="%s"/>' % (fx + math.cos(a*1.257)*0.07, fy + math.sin(a*1.257)*0.07, fc))
            o.append('<circle cx="%.2f" cy="%.2f" r="0.045" fill="#f9ca24"/>' % (fx, fy))
    o.append('</g>')
    return ''.join(o)


def boston_terriers():
    rng = random.Random(606)
    cols, rows = 5, 3
    tw, th = W/cols, H/rows
    bgs = ['#ffd6a5', '#caffbf', '#9bf6ff', '#bdb2ff', '#ffc6ff', '#fdffb6', '#a0c4ff', '#ffadad', '#b9fbc0', '#f1c0e8', '#cfbaf0', '#90dbf4', '#ffcfd2', '#98f5e1', '#fbf8cc']
    rng.shuffle(bgs)
    accs = ['bandana', 'bowtie', 'collar', 'party', 'flowers', None]
    coats = ['#1d1c22', '#1d1c22', '#1d1c22', '#26201f', '#2f2421']
    b = []
    defs = ''
    k = 0
    for r in range(rows):
        for c in range(cols):
            x0, y0 = c*tw, r*th
            bg = bgs[k]
            clip = 'tile%d' % k
            defs += '<clipPath id="%s"><rect x="%.0f" y="%.0f" width="%.0f" height="%.0f"/></clipPath>' % (clip, x0, y0, tw, th)
            b.append('<g clip-path="url(#%s)">' % clip)
            b.append('<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="%s"/>' % (x0, y0, tw, th, bg))
            pat = k % 4
            dark = shade(bg, -0.13)
            if pat == 0:
                for i in range(12):
                    for j in range(10):
                        b.append('<circle cx="%.0f" cy="%.0f" r="9" fill="%s"/>' % (x0 + 30 + i*55 + (j % 2)*27, y0 + 30 + j*70, dark))
            elif pat == 1:
                for i in range(-14, 14):
                    b.append('<path d="M%.0f %.0f l%.0f %.0f" stroke="%s" stroke-width="22"/>' % (x0 + i*60, y0, th, th, dark))
            elif pat == 2:
                for i in range(7):
                    for j in range(6):
                        px, py = x0 + 50 + i*95 + (j % 2)*45, y0 + 50 + j*115
                        b.append('<g fill="%s" transform="translate(%.0f %.0f) rotate(%d)"><ellipse cx="0" cy="8" rx="16" ry="13"/><circle cx="-15" cy="-10" r="6"/><circle cx="-5" cy="-17" r="6"/><circle cx="7" cy="-17" r="6"/><circle cx="16" cy="-9" r="6"/></g>' % (dark, px, py, rng.randint(-30, 30)))
            else:
                for i in range(8):
                    for j in range(8):
                        if (i + j) % 2 == 0:
                            b.append('<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="%s"/>' % (x0 + i*tw/8, y0 + j*th/8, tw/8, th/8, dark))
            # dog
            acc = accs[k % len(accs)]
            coat = rng.choice(coats)
            b.append('<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="%.0f" fill="#fff" opacity="0.35"/>' % (x0 + tw/2, y0 + th*0.48, tw*0.36, th*0.42))
            s = rng.uniform(170, 195)
            b.append(boston_head(x0 + tw/2 + rng.uniform(-20, 20), y0 + th*0.5 + 20, s, rng, coat, acc, rng.uniform(-10, 10), bg))
            if k % 3 == 1:
                bx, by = x0 + tw*0.8, y0 + th*0.82
                bc = rng.choice(['#e74c3c', '#2ecc71', '#f1c40f', '#3498db'])
                b.append('<circle cx="%.0f" cy="%.0f" r="48" fill="%s" stroke="%s" stroke-width="5"/><path d="M%.0f %.0f q48 30 96 0" stroke="#fff" stroke-width="7" fill="none"/>' % (bx, by, bc, shade(bc, -0.3), bx - 48, by - 5))
            elif k % 3 == 2:
                bx, by = x0 + tw*0.18, y0 + th*0.83
                b.append('<g transform="translate(%.0f %.0f) rotate(-25)"><rect x="-55" y="-14" width="110" height="28" rx="14" fill="#f5e6c8" stroke="#c8a96e" stroke-width="4"/><circle cx="-58" cy="-14" r="20" fill="#f5e6c8" stroke="#c8a96e" stroke-width="4"/><circle cx="-58" cy="14" r="20" fill="#f5e6c8" stroke="#c8a96e" stroke-width="4"/><circle cx="58" cy="-14" r="20" fill="#f5e6c8" stroke="#c8a96e" stroke-width="4"/><circle cx="58" cy="14" r="20" fill="#f5e6c8" stroke="#c8a96e" stroke-width="4"/><rect x="-55" y="-12" width="110" height="24" fill="#f5e6c8"/></g>' % (bx, by))
            b.append('</g>')
            b.append('<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="none" stroke="#fffaf0" stroke-width="16"/>' % (x0, y0, tw, th))
            k += 1
    return wrap(defs, ''.join(b), grain=0.05)


# ---------------------------------------------------------------- Flowers
def daisy(x, y, r, rng, col='#ffffff', center='#f6c90e'):
    o = []
    n = rng.randint(12, 18)
    for i in range(n):
        a = 360*i/n + rng.uniform(-4, 4)
        o.append('<ellipse cx="0" cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" stroke="%s" stroke-width="%.1f" transform="rotate(%.0f)"/>' % (-r*0.55, r*0.17, r*0.5, col, shade(col, -0.18), r*0.02, a))
    o.append('<circle r="%.1f" fill="%s"/><circle r="%.1f" fill="%s" opacity="0.6" cx="%.1f" cy="%.1f"/>' % (r*0.25, center, r*0.12, shade(center, -0.25), r*0.05, r*0.05))
    return '<g transform="translate(%.0f %.0f)">%s</g>' % (x, y, ''.join(o))


def tulip(x, y, r, col):
    d = shade(col, -0.25)
    return ('<g transform="translate(%.0f %.0f) scale(%.2f)"><path d="M-40 -10 C-45 40 45 40 40 -10 L28 -50 L14 -18 L0 -55 L-14 -18 L-28 -50Z" fill="%s" stroke="%s" stroke-width="2"/>'
            '<path d="M0 -55 C-20 -20 -18 20 0 30 C18 20 20 -20 0 -55Z" fill="%s" opacity="0.6"/></g>') % (x, y, r/40, col, d, shade(col, 0.15))


def poppy(x, y, r, rng, col='#e63022'):
    o = []
    for i in range(5):
        a = i*72 + rng.uniform(-10, 10)
        o.append('<ellipse cx="0" cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" stroke="%s" stroke-width="%.1f" transform="rotate(%.0f)"/>' % (-r*0.45, r*0.48, r*0.5, jitter_color(col, rng, 0.05), shade(col, -0.3), r*0.02, a))
    o.append('<circle r="%.1f" fill="#2a1a1a"/>' % (r*0.2))
    for i in range(10):
        a = i*36
        o.append('<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#1a1010"/>' % (math.cos(math.radians(a))*r*0.28, math.sin(math.radians(a))*r*0.28, r*0.04))
    return '<g transform="translate(%.0f %.0f)">%s</g>' % (x, y, ''.join(o))


def sunflower(x, y, r, rng):
    o = []
    for layer, (n, rr) in enumerate([(22, 1.0), (22, 0.85)]):
        for i in range(n):
            a = 360*i/n + layer*8
            c = '#f7b500' if layer == 0 else '#ffd23f'
            o.append('<ellipse cx="0" cy="%.1f" rx="%.1f" ry="%.1f" fill="%s" stroke="#c88a00" stroke-width="%.1f" transform="rotate(%.0f)"/>' % (-r*0.55*rr, r*0.13, r*0.35*rr, c, r*0.015, a))
    o.append('<circle r="%.1f" fill="#5a3410"/>' % (r*0.35))
    for i in range(60):
        a = i*137.5
        rr = r*0.33*math.sqrt(i/60)
        o.append('<circle cx="%.1f" cy="%.1f" r="%.1f" fill="#3a200a"/>' % (math.cos(math.radians(a))*rr, math.sin(math.radians(a))*rr, r*0.025))
    return '<g transform="translate(%.0f %.0f)">%s</g>' % (x, y, ''.join(o))


def lavender(x, y, h, rng):
    o = ['<path d="M%.0f %.0f q%.0f %.0f %.0f %.0f" stroke="#6b8e4e" stroke-width="5" fill="none"/>' % (x, y, rng.uniform(-10, 10), -h/2, rng.uniform(-15, 15), -h)]
    for i in range(14):
        t = 0.45 + i/14*0.55
        bx = x + rng.uniform(-3, 3)
        by = y - h*t
        o.append('<ellipse cx="%.0f" cy="%.0f" rx="9" ry="13" fill="%s"/>' % (bx + (8 if i % 2 else -8), by, jitter_color('#8e6fc8', rng, 0.08)))
    return ''.join(o)


def rose(x, y, r, col):
    d = shade(col, -0.3)
    o = ['<circle r="%.0f" fill="%s"/>' % (r, col)]
    for i in range(6):
        rr = r*(1 - i*0.15)
        o.append('<path d="M%.0f 0 A%.0f %.0f 0 0 1 %.0f %.0f" stroke="%s" stroke-width="%.1f" fill="none" transform="rotate(%d)"/>' % (-rr, rr, rr, rr*0.7, rr*0.6, d, r*0.06, i*67))
    return '<g transform="translate(%.0f %.0f)">%s</g>' % (x, y, ''.join(o))


def butterfly(x, y, s, rot, c1, c2):
    return ('<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)">'
            '<path d="M0 0 C-30 -60 -80 -40 -60 0 C-80 30 -30 50 0 10Z" fill="%s" stroke="#222" stroke-width="3"/>'
            '<path d="M0 0 C30 -60 80 -40 60 0 C80 30 30 50 0 10Z" fill="%s" stroke="#222" stroke-width="3"/>'
            '<circle cx="-38" cy="-18" r="9" fill="%s"/><circle cx="38" cy="-18" r="9" fill="%s"/>'
            '<ellipse rx="5" ry="26" fill="#222"/><path d="M-2 -24 q-10 -20 -18 -22 M2 -24 q10 -20 18 -22" stroke="#222" stroke-width="2" fill="none"/></g>') % (x, y, rot, s, c1, c1, c2, c2)


def flower_garden():
    rng = random.Random(707)
    defs = '''<linearGradient id="sky" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#7ec8f2"/><stop offset="1" stop-color="#d9f0fb"/></linearGradient>'''
    b = ['<rect width="%d" height="%d" fill="url(#sky)"/>' % (W, H)]
    b.append('<circle cx="2700" cy="220" r="120" fill="#fff6c9" filter="url(#glow)"/>')
    for i in range(8):
        b.append('<path d="%s" fill="#ffffff" opacity="0.9"/>' % blob(rng.uniform(0, W), rng.uniform(80, 420), rng.uniform(70, 120), rng, 10, 0.25, 2.4, 0.8))
    for col, base, amp in [('#9fcf8a', 640, 70), ('#7fbf6a', 720, 60), ('#66ad55', 800, 40)]:
        p = [(0, H)] + [(i, base - amp*math.sin(i/500 + base) - rng.uniform(0, 10)) for i in range(0, W+101, 100)] + [(W, H)]
        b.append('<polygon points="%s" fill="%s"/>' % (pts(p), col))
    # farmhouse + trees in distance
    b.append('<rect x="560" y="610" width="140" height="90" fill="#c0392b"/><polygon points="550,612 630,560 710,612" fill="#7f2a20"/><rect x="615" y="650" width="26" height="50" fill="#5b2a1a"/>')
    for tx in (300, 380, 2200, 2280, 2350):
        b.append('<rect x="%d" y="640" width="12" height="60" fill="#5a3d26"/><circle cx="%d" cy="620" r="45" fill="#4f8f3f"/>' % (tx, tx + 6))
    b.append('<rect y="780" width="3200" height="1220" fill="#5a9a46"/>')
    # rows of flowers from back to front
    items = []
    for i in range(1400):
        y = rng.uniform(760, 2050)
        items.append((y, rng.uniform(-40, W+40)))
    items.sort()
    bands = [(0, 1000, 'lavender'), (1000, 1200, 'tulip'), (1200, 1450, 'mix'), (1450, 2100, 'mix2')]
    tulipc = ['#e84a5f', '#ff847c', '#fecea8', '#f9d423', '#ff4e50', '#c06c84', '#f8b195']
    for y, x in items:
        t = (y - 760)/1300
        sc = 0.35 + t*1.4
        # stem + leaves
        stem_h = 70*sc
        if y < 1000:
            b.append(lavender(x, y, 110*sc, rng))
            continue
        b.append('<path d="M%.0f %.0f l%.0f %.0f" stroke="#3e7d32" stroke-width="%.1f"/>' % (x, y, rng.uniform(-6, 6), -stem_h, max(2, 4*sc)))
        if rng.random() < 0.5:
            b.append(leaf_s(x, y - stem_h*0.4, 30*sc, rng.uniform(-70, -20) if rng.random() < 0.5 else rng.uniform(20, 70)))
        fy = y - stem_h
        if y < 1200:
            b.append(tulip(x, fy, 30*sc, rng.choice(tulipc)))
        else:
            k = rng.random()
            if k < 0.25: b.append(daisy(x, fy, 34*sc, rng, rng.choice(['#ffffff', '#ffffff', '#ffe4f2', '#fff5c0'])))
            elif k < 0.45: b.append(poppy(x, fy, 34*sc, rng, rng.choice(['#e63022', '#ff6a3d', '#ff3d6e'])))
            elif k < 0.58: b.append(sunflower(x, fy - 20*sc, 42*sc, rng))
            elif k < 0.75: b.append(rose(x, fy, 22*sc, rng.choice(['#e84a5f', '#ff8fab', '#ffb3c6', '#c9184a', '#fff0f3'])))
            elif k < 0.88: b.append(tulip(x, fy, 30*sc, rng.choice(tulipc)))
            else: b.append(daisy(x, fy, 30*sc, rng, rng.choice(['#b388eb', '#8093f1', '#72ddf7']), '#ffe066'))
    for i in range(14):
        c1 = rng.choice(['#ff9f1c', '#2ec4b6', '#e71d36', '#9b5de5', '#00bbf9', '#f15bb5'])
        b.append(butterfly(rng.uniform(100, W-100), rng.uniform(850, 1800), rng.uniform(0.8, 1.4), rng.uniform(-30, 30), c1, '#fff'))
    return wrap(defs, ''.join(b), grain=0.05)


def leaf_s(x, y, s, rot):
    return '<path d="M0 0 C-0.4 -0.3 -0.4 -0.8 0 -1.2 C0.4 -0.8 0.4 -0.3 0 0Z" fill="#4c9a3a" transform="translate(%.0f %.0f) rotate(%.0f) scale(%.1f)"/>' % (x, y, rot, s)


# ---------------------------------------------------------------- Ocean
def fish(x, y, s, rng, flip=False):
    kind = rng.choice(['clown', 'tang', 'stripe', 'angel', 'plain'])
    sx = -s if flip else s
    o = ['<g transform="translate(%.0f %.0f) scale(%.2f %.2f)">' % (x, y, sx, s)]
    if kind == 'clown':
        o.append('<path d="M-60 0 C-40 -38 30 -40 55 0 C30 40 -40 38 -60 0Z" fill="#ff7b1c" stroke="#222" stroke-width="3"/><path d="M-60 0 L-90 -28 L-84 0 L-90 28Z" fill="#ff7b1c" stroke="#222" stroke-width="3"/>')
        for bx in (-30, 0, 30):
            o.append('<path d="M%d -34 C%d -10 %d 10 %d 34 L%d 34 C%d 10 %d -10 %d -34Z" fill="#fff" stroke="#222" stroke-width="3"/>' % (bx - 6, bx - 10, bx - 10, bx - 6, bx + 8, bx + 4, bx + 4, bx + 8))
    elif kind == 'tang':
        o.append('<path d="M-60 0 C-40 -45 30 -45 55 0 C30 45 -40 45 -60 0Z" fill="#1f4fd8" stroke="#0a1f5a" stroke-width="3"/><path d="M-60 0 L-92 -30 L-86 0 L-92 30Z" fill="#ffd60a"/><path d="M-40 -10 C-10 -30 20 -20 30 0 C10 -10 -20 0 -40 10Z" fill="#0a1f5a"/>')
    elif kind == 'stripe':
        c = rng.choice(['#ffd166', '#06d6a0', '#ef476f'])
        o.append('<path d="M-60 0 C-40 -32 30 -34 55 0 C30 34 -40 32 -60 0Z" fill="%s" stroke="#333" stroke-width="3"/><path d="M-60 0 L-88 -22 L-82 0 L-88 22Z" fill="%s"/>' % (c, c))
        for k in range(4):
            o.append('<path d="M%d -28 v56" stroke="#2b2d42" stroke-width="7" opacity="0.6"/>' % (-35 + k*20))
    elif kind == 'angel':
        o.append('<path d="M-40 0 C-30 -40 30 -50 45 0 C30 50 -30 40 -40 0Z" fill="#ffe066" stroke="#333" stroke-width="3"/><path d="M-10 -38 L10 -95 L25 -35Z M-10 38 L10 95 L25 35Z" fill="#ffe066" stroke="#333" stroke-width="3"/><path d="M-40 0 L-65 -20 L-65 20Z" fill="#4cc9f0"/><path d="M-20 -32 v64 M5 -42 v84" stroke="#4cc9f0" stroke-width="8"/>')
    else:
        c = rng.choice(['#f15bb5', '#9b5de5', '#00bbf9', '#fee440', '#fb5607'])
        o.append('<path d="M-60 0 C-40 -30 30 -32 55 0 C30 32 -40 30 -60 0Z" fill="%s" stroke="%s" stroke-width="3"/><path d="M-60 0 L-86 -22 L-80 0 L-86 22Z" fill="%s"/>' % (c, shade(c, -0.4), shade(c, -0.15)))
        o.append('<path d="M-30 18 C0 28 30 18 45 5" stroke="#fff" stroke-width="5" fill="none" opacity="0.5"/>')
    o.append('<circle cx="35" cy="-6" r="9" fill="#fff"/><circle cx="37" cy="-6" r="5" fill="#111"/></g>')
    return ''.join(o)


def coral_branch(x, y, h, col, rng, depth=5):
    out = []

    def br(x, y, l, a, w, d):
        if d == 0:
            out.append('<circle cx="%.0f" cy="%.0f" r="%.1f" fill="%s"/>' % (x, y, w*0.7, shade(col, 0.2)))
            return
        x2, y2 = x + math.cos(a)*l, y + math.sin(a)*l
        out.append('<path d="M%.0f %.0f L%.0f %.0f" stroke="%s" stroke-width="%.1f" stroke-linecap="round"/>' % (x, y, x2, y2, col, w))
        for s in (-1, 1):
            br(x2, y2, l*0.75, a + s*rng.uniform(0.25, 0.6), w*0.8, d-1)
    br(x, y, h*0.35, -math.pi/2 + rng.uniform(-0.2, 0.2), h*0.08, depth)
    return ''.join(out)


def ocean():
    rng = random.Random(808)
    defs = '''<linearGradient id="sea" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#4cc9f0"/><stop offset="0.5" stop-color="#1b75bb"/><stop offset="1" stop-color="#0b3c74"/></linearGradient>
<linearGradient id="ray" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#ffffff" stop-opacity="0.5"/><stop offset="1" stop-color="#ffffff" stop-opacity="0"/></linearGradient>
<linearGradient id="sand" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#e8cf98"/><stop offset="1" stop-color="#c9a86a"/></linearGradient>
<radialGradient id="jelly" cx="0.5" cy="0.3"><stop offset="0" stop-color="#ffd6f5" stop-opacity="0.95"/><stop offset="1" stop-color="#d16ba5" stop-opacity="0.6"/></radialGradient>'''
    b = ['<rect width="%d" height="%d" fill="url(#sea)"/>' % (W, H)]
    for i in range(9):
        x = rng.uniform(0, W)
        b.append('<polygon points="%.0f,0 %.0f,0 %.0f,1500 %.0f,1500" fill="url(#ray)" opacity="0.5"/>' % (x, x + rng.uniform(60, 160), x + 400, x + 150))
    # surface ripples
    for i in range(40):
        b.append('<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="8" fill="#bdf0ff" opacity="0.35"/>' % (rng.uniform(0, W), rng.uniform(5, 80), rng.uniform(40, 140)))
    # distant rocks
    for i in range(10):
        b.append('<path d="%s" fill="#145a8c" opacity="0.7"/>' % blob(rng.uniform(0, W), 1450, rng.uniform(180, 320), rng, 9, 0.25, 1.4, 0.8))
    # sand
    b.append('<path d="M0 1500 Q800 1430 1600 1490 T3200 1470 L3200 2000 L0 2000Z" fill="url(#sand)"/>')
    for i in range(300):
        b.append('<circle cx="%.0f" cy="%.0f" r="%.1f" fill="#a8864e" opacity="0.5"/>' % (rng.uniform(0, W), rng.uniform(1520, 2000), rng.uniform(2, 6)))
    # seaweed
    for i in range(30):
        x = rng.uniform(0, W)
        h = rng.uniform(300, 700)
        p = []
        for k in range(10):
            t = k/9
            p.append((x + math.sin(t*6 + i)*30*t, 1560 - h*t))
        c = rng.choice(['#2a9d8f', '#3a7d44', '#6a994e', '#1b7a5a'])
        b.append('<path d="%s" stroke="%s" stroke-width="%.0f" fill="none" stroke-linecap="round"/>' % (smooth_path(p, False), c, rng.uniform(14, 26)))
    # corals
    coral_cols = ['#ff6b6b', '#ff9f43', '#f368e0', '#feca57', '#ff7eb3', '#e056fd', '#ff5e57']
    for i in range(22):
        x = rng.uniform(0, W)
        y = rng.uniform(1560, 1950)
        h = rng.uniform(250, 520)
        col = rng.choice(coral_cols)
        k = rng.random()
        if k < 0.4:
            b.append(coral_branch(x, y, h, col, rng))
        elif k < 0.65:
            # brain coral
            b.append('<path d="%s" fill="%s" stroke="%s" stroke-width="6"/>' % (blob(x, y - h*0.15, h*0.28, rng, 12, 0.1, 1.3, 0.8), col, shade(col, -0.3)))
            for j in range(6):
                b.append('<path d="M%.0f %.0f q30 -25 60 0 t60 0" stroke="%s" stroke-width="5" fill="none"/>' % (x - h*0.3 + rng.uniform(0, h*0.2), y - h*0.25 + j*h*0.06, shade(col, -0.25)))
        elif k < 0.85:
            # fan coral
            for j in range(18):
                a = -math.pi + j/17*math.pi
                b.append('<path d="M%.0f %.0f Q%.0f %.0f %.0f %.0f" stroke="%s" stroke-width="6" fill="none" opacity="0.9"/>' % (x, y, x + math.cos(a)*h*0.3, y + math.sin(a)*h*0.5, x + math.cos(a)*h*0.5, y + math.sin(a)*h*0.75, col))
        else:
            # anemone
            for j in range(26):
                a = -math.pi + j/25*math.pi
                b.append('<path d="M%.0f %.0f q%.0f %.0f %.0f %.0f" stroke="%s" stroke-width="12" fill="none" stroke-linecap="round"/>' % (x, y, math.cos(a)*40, math.sin(a)*90, math.cos(a)*h*0.3, math.sin(a)*h*0.4, col))
            b.append('<ellipse cx="%.0f" cy="%.0f" rx="%.0f" ry="40" fill="%s"/>' % (x, y, h*0.25, shade(col, -0.3)))
    # starfish and shells
    for i in range(8):
        x, y = rng.uniform(0, W), rng.uniform(1700, 1980)
        c = rng.choice(['#ff7f50', '#ff6b81', '#ffa502', '#e84393'])
        p = []
        for k in range(10):
            a = math.radians(k*36 - 90)
            r = 60 if k % 2 == 0 else 26
            p.append((math.cos(a)*r, math.sin(a)*r))
        b.append('<g transform="translate(%.0f %.0f) rotate(%.0f)"><path d="%s" fill="%s" stroke="%s" stroke-width="4"/></g>' % (x, y, rng.uniform(0, 72), smooth_path(p), c, shade(c, -0.3)))
    # turtle
    tx, ty = 2200, 650
    b.append('<g transform="translate(%d %d) rotate(-12)"><ellipse cx="-150" cy="60" rx="80" ry="34" fill="#7fb069" transform="rotate(30 -150 60)"/><ellipse cx="150" cy="60" rx="80" ry="30" fill="#7fb069" transform="rotate(-30 150 60)"/>'
             '<ellipse cx="-120" cy="-70" rx="110" ry="40" fill="#7fb069" transform="rotate(-35 -120 -70)"/><ellipse cx="120" cy="-70" rx="110" ry="40" fill="#7fb069" transform="rotate(35 120 -70)"/>'
             '<ellipse cx="0" cy="-170" rx="55" ry="60" fill="#8cbf73"/><circle cx="-22" cy="-185" r="9" fill="#111"/><circle cx="22" cy="-185" r="9" fill="#111"/>'
             '<ellipse cx="0" cy="0" rx="170" ry="150" fill="#6b4f2a" stroke="#3e2c14" stroke-width="8"/>' % (tx, ty))
    for (hx, hy, hr) in [(0, 0, 50), (-90, -40, 40), (90, -40, 40), (-90, 55, 40), (90, 55, 40), (0, -100, 38), (0, 100, 38)]:
        b.append('<polygon points="%s" fill="#a67c3d" stroke="#3e2c14" stroke-width="5"/>' % pts([(hx + math.cos(math.radians(a))*hr, hy + math.sin(math.radians(a))*hr) for a in range(0, 360, 60)]))
    b.append('</g>')
    # jellyfish
    for (jx, jy, js) in [(500, 420, 1.2), (780, 300, 0.8), (2900, 300, 1.0)]:
        o = '<g transform="translate(%d %d) scale(%.1f)" opacity="0.9">' % (jx, jy, js)
        for k in range(7):
            o += '<path d="M%d 0 q20 60 -5 120 t5 120" stroke="#f7b2e0" stroke-width="6" fill="none" opacity="0.8"/>' % (-75 + k*25)
        o += '<path d="M-100 0 C-100 -130 100 -130 100 0 Q75 18 50 0 Q25 18 0 0 Q-25 18 -50 0 Q-75 18 -100 0Z" fill="url(#jelly)"/></g>'
        b.append(o)
    # fish schools
    for i in range(55):
        x, y = rng.uniform(50, W-50), rng.uniform(250, 1500)
        b.append(fish(x, y, rng.uniform(0.6, 1.6), rng, rng.random() < 0.5))
    for i in range(70):
        x, y = rng.uniform(0, W), rng.uniform(100, 1800)
        r = rng.uniform(5, 18)
        b.append('<circle cx="%.0f" cy="%.0f" r="%.0f" fill="none" stroke="#e0f7ff" stroke-width="3" opacity="0.7"/>' % (x, y, r))
    return wrap(defs, ''.join(b), grain=0.05)


# ---------------------------------------------------------------- Forest
def pine(x, base, h, col, rng, trunk=None):
    w = h*0.34
    o = []
    o.append('<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" fill="%s"/>' % (x - h*0.02, base - h*0.14, h*0.04, h*0.14, trunk or shade(col, -0.35)))
    tiers = 5
    for i in range(tiers):
        t0 = i/tiers
        yb = base - h*0.1 - h*0.72*t0
        yt = yb - h*0.34 if i < tiers - 1 else base - h
        ww = w*(1 - t0*0.78)*rng.uniform(0.92, 1.05)
        c = shade(col, 0.06*i/tiers)
        o.append('<path d="M%.0f %.0f Q%.0f %.0f %.0f %.0f Q%.0f %.0f %.0f %.0f Q%.0f %.0f %.0f %.0fZ" fill="%s"/>' % (
            x, yt, x - ww*0.35, yb - h*0.12, x - ww, yb, x, yb - h*0.05, x + ww, yb, x + ww*0.35, yb - h*0.12, x, yt, c))
        o.append('<path d="M%.0f %.0f Q%.0f %.0f %.0f %.0f Q%.0f %.0f %.0f %.0f" fill="%s" opacity="0.35"/>' % (
            x, yt, x + ww*0.35, yb - h*0.12, x + ww, yb, x + ww*0.4, yb - h*0.03, x, yb - h*0.05, shade(col, -0.25)))
    return ''.join(o)


def mushroom(x, y, s, col):
    return ('<g transform="translate(%.0f %.0f) scale(%.2f)"><path d="M-12 0 L-9 -50 L9 -50 L12 0Z" fill="#f3ead8" stroke="#a89880" stroke-width="2"/>'
            '<path d="M-48 -45 C-48 -95 48 -95 48 -45 Z" fill="%s" stroke="%s" stroke-width="3"/>'
            '<circle cx="-20" cy="-65" r="8" fill="#fff"/><circle cx="12" cy="-75" r="6" fill="#fff"/><circle cx="25" cy="-55" r="7" fill="#fff"/></g>') % (x, y, s, col, shade(col, -0.3))


def fern(x, y, s, rot, col):
    o = ['<g transform="translate(%.0f %.0f) rotate(%.0f) scale(%.2f)"><path d="M0 0 Q10 -100 0 -200" stroke="%s" stroke-width="4" fill="none"/>' % (x, y, rot, s, col)]
    for i in range(14):
        t = i/14
        yy = -t*190
        l = 70*(1 - t)
        o.append('<path d="M0 %.0f q%.0f -10 %.0f -%.0f" stroke="%s" stroke-width="9" fill="none" stroke-linecap="round"/>' % (yy, -l*0.5, -l, l*0.3, col))
        o.append('<path d="M0 %.0f q%.0f -10 %.0f -%.0f" stroke="%s" stroke-width="9" fill="none" stroke-linecap="round"/>' % (yy, l*0.5, l, l*0.3, col))
    o.append('</g>')
    return ''.join(o)


def forest():
    rng = random.Random(909)
    defs = '''<linearGradient id="sky" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#ffd9a0"/><stop offset="0.5" stop-color="#fbe8c8"/><stop offset="1" stop-color="#d7e8d0"/></linearGradient>
<linearGradient id="ray" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#fff6d0" stop-opacity="0.65"/><stop offset="1" stop-color="#fff6d0" stop-opacity="0"/></linearGradient>
<linearGradient id="mist" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#f4f1e0" stop-opacity="0"/><stop offset="1" stop-color="#f4f1e0" stop-opacity="0.7"/></linearGradient>
<linearGradient id="water" x1="0" y1="0" x2="0" y2="1"><stop offset="0" stop-color="#8fc6d8"/><stop offset="1" stop-color="#4c8fa8"/></linearGradient>'''
    b = ['<rect width="%d" height="%d" fill="url(#sky)"/>' % (W, H)]
    b.append('<circle cx="2100" cy="520" r="160" fill="#fff8e0" filter="url(#glow)"/>')
    # mountains
    for col, base, amp in [('#b8c8c0', 900, 300), ('#9db5aa', 1000, 220)]:
        p = [(0, H)]
        x = 0
        while x <= W + 200:
            p.append((x, base - rng.uniform(0.3, 1)*amp))
            x += rng.uniform(150, 300)
        p.append((W, H))
        b.append('<polygon points="%s" fill="%s"/>' % (pts(p), col))
    layers = [('#89a898', 1150, 160, 260, 18), ('#6f9282', 1250, 220, 340, 15), ('#557a69', 1360, 300, 440, 12), ('#3c6150', 1480, 420, 600, 9), ('#284a3b', 1650, 650, 950, 5)]
    for li, (col, base, hmin, hmax, n) in enumerate(layers):
        b.append('<rect y="%d" width="3200" height="%d" fill="%s"/>' % (base - 10, H - base + 10, col))
        xs = sorted(rng.uniform(-100, W+100) for _ in range(n*3))
        for x in xs:
            b.append(pine(x, base + rng.uniform(-20, 20), rng.uniform(hmin, hmax), jitter_color(col, rng, 0.04), rng))
        if li < 4:
            b.append('<rect y="%d" width="3200" height="260" fill="url(#mist)"/>' % (base - 200))
        if li == 2:
            for i in range(8):
                x = 1750 + i*90
                b.append('<polygon points="%.0f,0 %.0f,0 %.0f,2000 %.0f,2000" fill="url(#ray)" opacity="0.5"/>' % (x, x + 40, x - 600 + i*40, x - 750 + i*40))
    # stream
    b.append('<path d="M1300 1660 C1500 1700 1200 1800 1400 2000 L1900 2000 C1650 1820 1800 1720 1500 1660Z" fill="url(#water)"/>')
    for i in range(14):
        y = 1690 + i*22
        b.append('<path d="M%.0f %.0f h%.0f" stroke="#e5f6fb" stroke-width="4" opacity="0.6"/>' % (1420 + rng.uniform(-60, 120) + i*12, y, rng.uniform(40, 120)))
    # rocks
    for i in range(10):
        x, y = rng.uniform(1100, 2100), rng.uniform(1700, 2000)
        b.append('<path d="%s" fill="%s" stroke="#4a5550" stroke-width="4"/>' % (blob(x, y, rng.uniform(30, 70), rng, 8, 0.2, 1.4, 0.8), jitter_color('#8a948f', rng)))
    # deer silhouette
    dx, dy = 2450, 1620
    b.append('<g transform="translate(%d %d)" fill="#1e3a2e"><ellipse cx="0" cy="-120" rx="110" ry="50"/><rect x="-90" y="-100" width="16" height="110"/><rect x="-60" y="-100" width="16" height="110"/><rect x="60" y="-100" width="16" height="110"/><rect x="85" y="-100" width="16" height="110"/>'
             '<path d="M70 -150 L110 -260 L150 -250 L120 -140Z"/><ellipse cx="140" cy="-265" rx="40" ry="24"/><path d="M120 -285 l-20 -60 l-20 10 m20 -10 l10 -40 M150 -285 l20 -60 l20 10 m-20 -10 l-10 -40" stroke="#1e3a2e" stroke-width="9" fill="none"/>'
             '<path d="M110 -280 l-25 -15 l10 20Z"/></g>' % (dx, dy))
    # foreground ferns, mushrooms
    for i in range(30):
        x = rng.uniform(0, W)
        if 1200 < x < 2000:
            continue
        b.append(fern(x, 2010, rng.uniform(1.0, 1.8), rng.uniform(-30, 30), rng.choice(['#3f7d3a', '#4f8f45', '#2f6a30'])))
    for i in range(10):
        x = rng.uniform(100, W-100)
        if 1250 < x < 1950:
            continue
        b.append(mushroom(x, rng.uniform(1880, 1990), rng.uniform(0.9, 1.5), rng.choice(['#d63031', '#e17055', '#c0392b'])))
    # big foreground trunks
    # fallen log + wildflowers
    b.append('<g transform="translate(420 1880) rotate(-6)"><rect x="-260" y="-55" width="520" height="110" rx="55" fill="#6b4a32" stroke="#3b2a1e" stroke-width="6"/><ellipse cx="260" cy="0" rx="40" ry="55" fill="#c8a27a" stroke="#3b2a1e" stroke-width="6"/><ellipse cx="260" cy="0" rx="22" ry="32" fill="none" stroke="#8a6a4a" stroke-width="5"/><path d="M-200 -20 h300 M-150 15 h280" stroke="#4a3322" stroke-width="6"/></g>')
    for i in range(70):
        x = rng.uniform(0, W)
        if 1250 < x < 1950:
            continue
        y = rng.uniform(1700, 1990)
        c = rng.choice(['#ffffff', '#fff3a0', '#f7a1c4', '#b8a1f7', '#ffd166'])
        b.append('<path d="M%.0f %.0f v-%.0f" stroke="#3f7d3a" stroke-width="4"/>' % (x, y, 30))
        for a in range(5):
            b.append('<circle cx="%.0f" cy="%.0f" r="7" fill="%s"/>' % (x + math.cos(a*1.257)*8, y - 30 + math.sin(a*1.257)*8, c))
        b.append('<circle cx="%.0f" cy="%.0f" r="5" fill="#f4a300"/>' % (x, y - 30))
    for i in range(12):
        bx, by = rng.uniform(300, 2900), rng.uniform(250, 900)
        b.append('<path d="M%.0f %.0f q12 -12 24 0 q12 -12 24 0" stroke="#3b3b3b" stroke-width="4" fill="none"/>' % (bx, by))
    return wrap(defs, ''.join(b), grain=0.05)


# ---------------------------------------------------------------- Stained glass
def stained_glass():
    from scipy.spatial import Voronoi
    import numpy as np
    rng = random.Random(1010)
    b = []
    defs = ''
    # stone wall background
    b.append('<rect width="%d" height="%d" fill="#8a8274"/>' % (W, H))
    for row in range(0, H, 110):
        off = (row//110 % 2)*120
        for col in range(-240, W, 240):
            c = jitter_color('#9d9585', rng, 0.06)
            b.append('<rect x="%d" y="%d" width="232" height="102" rx="8" fill="%s"/>' % (col + off, row, c))
            for k in range(3):
                b.append('<circle cx="%.0f" cy="%.0f" r="%.0f" fill="%s" opacity="0.35"/>' % (col + off + rng.uniform(10, 220), row + rng.uniform(10, 92), rng.uniform(3, 10), shade(c, -0.2)))
    # window shape (gothic arch)
    cx = W/2
    wl, wr, wt, wb = 760, 2440, 120, 1800
    spring = 820
    arch ='M%d %d L%d %d Q%d %d %d %d Q%d %d %d %d L%d %dZ' % (wl, wb, wl, spring, wl, wt + 150, cx, wt, wr, wt + 150, wr, spring, wr, wb)
    defs += '<clipPath id="win"><path d="%s"/></clipPath>' % arch
    # frame
    b.append('<path d="%s" fill="none" stroke="#6b6355" stroke-width="120"/>' % arch)
    b.append('<path d="%s" fill="none" stroke="#c9bfa8" stroke-width="60"/>' % arch)
    # voronoi glass
    pts_ = []
    for i in range(520):
        pts_.append((rng.uniform(wl - 100, wr + 100), rng.uniform(wt - 100, wb + 100)))
    for x in range(wl - 300, wr + 400, 300):
        pts_ += [(x, wt - 400), (x, wb + 400)]
    for y in range(wt - 300, wb + 400, 300):
        pts_ += [(wl - 400, y), (wr + 400, y)]
    vor = Voronoi(np.array(pts_))
    sun = (cx, 760)
    glass = ['<g clip-path="url(#win)">']
    for i, reg_i in enumerate(vor.point_region):
        reg = vor.regions[reg_i]
        if not reg or -1 in reg:
            continue
        poly = [tuple(vor.vertices[j]) for j in reg]
        px, py = pts_[i]
        dist = math.hypot(px - sun[0], py - sun[1])
        ang = math.atan2(py - sun[1], px - sun[0])
        if dist < 260:
            c = rng.choice(['#fff4c2', '#ffe680', '#fff9e0', '#ffd34d'])
        elif dist < 520:
            c = rng.choice(['#ffc845', '#ffb02e', '#ffd966', '#ff9f1c'])
        else:
            ray = (int((ang + math.pi)/(2*math.pi)*16)) % 2
            if ray == 0:
                c = rng.choice(['#1e5fb4', '#2d7dd2', '#164a8c', '#3c8fe0', '#0f3d7a'])
            else:
                c = rng.choice(['#7b2cbf', '#5a189a', '#9d4edd', '#c2185b', '#2a9d8f', '#1b998b'])
            if py > 1450:
                c = rng.choice(['#2e7d32', '#43a047', '#1b5e20', '#66bb6a', '#7cb342'])
        g = 'sg%d' % i
        defs += '<radialGradient id="%s" cx="0.35" cy="0.3"><stop offset="0" stop-color="%s"/><stop offset="1" stop-color="%s"/></radialGradient>' % (g, shade(c, 0.25), shade(c, -0.12))
        glass.append('<polygon points="%s" fill="url(#%s)" stroke="#2a2522" stroke-width="10" stroke-linejoin="round"/>' % (pts(poly), g))
    # cross
    cw = 90
    cross = 'M%.0f %.0f h%d v%d h%d v%d h-%d v%d h-%d v-%d h-%d v-%d h%d Z' % (cx - cw/2, 380, cw, 260, 230, cw, 230, 760, cw, 760, 230, cw, 230)
    glass.append('<path d="%s" fill="#fff8e1" stroke="#2a2522" stroke-width="14"/>' % cross)
    for y in range(470, 1460, 110):
        glass.append('<path d="M%.0f %d h%d" stroke="#2a2522" stroke-width="8"/>' % (cx - cw/2, y, cw))
    for x in (cx - 230 + 100, cx + 130):
        glass.append('<path d="M%.0f %d v%d" stroke="#2a2522" stroke-width="8"/>' % (x, 640, cw))
    # dove
    glass.append('<g transform="translate(%.0f 300) scale(1.3)"><path d="M0 0 C-40 -10 -120 -60 -170 -40 C-120 -20 -80 10 -40 20 C-60 40 -40 60 0 50 C40 60 60 40 40 20 C80 10 120 -20 170 -40 C120 -60 40 -10 0 0Z" fill="#ffffff" stroke="#2a2522" stroke-width="8"/>'
                 '<circle cx="0" cy="-8" r="16" fill="#ffffff" stroke="#2a2522" stroke-width="6"/><path d="M-60 -25 L-90 -30 M60 -25 L90 -30 M-30 10 L-60 5 M30 10 L60 5" stroke="#2a2522" stroke-width="5"/></g>' % cx)
    # lilies
    for sx in (-1, 1):
        lx = cx + sx*560
        glass.append('<path d="M%.0f 1800 Q%.0f 1500 %.0f 1250" stroke="#2a2522" stroke-width="26" fill="none"/><path d="M%.0f 1800 Q%.0f 1500 %.0f 1250" stroke="#4caf50" stroke-width="14" fill="none"/>' % (lx, lx + sx*30, lx, lx, lx + sx*30, lx))
        for k, a in enumerate((-50, 0, 50)):
            glass.append('<g transform="translate(%.0f 1250) rotate(%d)"><path d="M0 0 C-40 -40 -30 -110 0 -150 C30 -110 40 -40 0 0Z" fill="#fffdf5" stroke="#2a2522" stroke-width="9"/></g>' % (lx, a))
        glass.append('<circle cx="%.0f" cy="1245" r="16" fill="#ffc107" stroke="#2a2522" stroke-width="6"/>' % lx)
    glass.append('</g>')
    b += glass
    # mullion lines
    b.append('<path d="%s" fill="none" stroke="#2a2522" stroke-width="20"/>' % arch)
    # sill + candles
    b.append('<rect x="620" y="1820" width="1960" height="90" fill="#bdb29a" stroke="#6b6355" stroke-width="8"/>')
    for x in (700, 860, 2340, 2500):
        h = rng.uniform(160, 240)
        b.append('<rect x="%d" y="%.0f" width="60" height="%.0f" fill="#f6efe0" stroke="#c9bfa8" stroke-width="4"/>' % (x, 1820 - h, h))
        b.append('<ellipse cx="%d" cy="%.0f" rx="14" ry="30" fill="#ffcf4a" filter="url(#glow)"/>' % (x + 30, 1820 - h - 34))
    # banner
    b.append('<path d="M900 1930 L2300 1930 L2350 1975 L2300 2000 L900 2000 L850 1975Z" fill="#7a1f2b"/>')
    b.append('<text x="%d" y="1986" font-family="Lora" font-weight="bold" font-size="54" fill="#f7e3b0" text-anchor="middle" letter-spacing="6">I AM THE LIGHT OF THE WORLD · JOHN 8:12</text>' % cx)
    return wrap(defs, ''.join(b), grain=0.05)


SCENES = [
    ('boston-terriers', 'Boston Terrier Pals', boston_terriers),
    ('flower-garden', 'Spring Flower Garden', flower_garden),
    ('ocean-reef', 'Coral Reef Friends', ocean),
    ('forest-morning', 'Misty Forest Morning', forest),
    ('stained-glass', 'Light of the World', stained_glass),
]
