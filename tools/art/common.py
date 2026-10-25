"""Shared helpers for generating Snugfit artwork as SVG and rendering it with Chromium."""
import math, random, colorsys

W, H = 3200, 2000


def hexrgb(h):
    h = h.lstrip('#')
    return tuple(int(h[i:i+2], 16) for i in (0, 2, 4))


def rgbhex(r, g, b):
    return '#%02x%02x%02x' % (max(0, min(255, int(r))), max(0, min(255, int(g))), max(0, min(255, int(b))))


def mix(a, b, t):
    ra, ga, ba = hexrgb(a); rb, gb, bb = hexrgb(b)
    return rgbhex(ra + (rb-ra)*t, ga + (gb-ga)*t, ba + (bb-ba)*t)


def shade(c, f):
    """f<0 darker, f>0 lighter"""
    return mix(c, '#000000', -f) if f < 0 else mix(c, '#ffffff', f)


def jitter_color(c, rng, amt=0.08):
    r, g, b = [x/255 for x in hexrgb(c)]
    h, l, s = colorsys.rgb_to_hls(r, g, b)
    h = (h + rng.uniform(-amt, amt)*0.3) % 1
    l = max(0, min(1, l + rng.uniform(-amt, amt)))
    s = max(0, min(1, s + rng.uniform(-amt, amt)))
    r, g, b = colorsys.hls_to_rgb(h, l, s)
    return rgbhex(r*255, g*255, b*255)


def pts(points):
    return ' '.join('%.1f,%.1f' % p for p in points)


def smooth_path(points, closed=True):
    """Catmull-Rom spline through points -> cubic bezier path"""
    n = len(points)
    if n < 3:
        return 'M' + ' L'.join('%.1f %.1f' % p for p in points)
    d = 'M%.1f %.1f' % points[0]
    rng_ = range(n) if closed else range(n-1)
    for i in rng_:
        p0 = points[(i-1) % n] if closed or i > 0 else points[0]
        p1 = points[i]
        p2 = points[(i+1) % n]
        p3 = points[(i+2) % n] if closed or i+2 < n else points[-1]
        c1 = (p1[0] + (p2[0]-p0[0])/6, p1[1] + (p2[1]-p0[1])/6)
        c2 = (p2[0] - (p3[0]-p1[0])/6, p2[1] - (p3[1]-p1[1])/6)
        d += ' C%.1f %.1f %.1f %.1f %.1f %.1f' % (c1 + c2 + p2)
    if closed:
        d += 'Z'
    return d


def blob(cx, cy, r, rng, n=10, var=0.25, sx=1.0, sy=1.0):
    p = []
    for i in range(n):
        a = 2*math.pi*i/n
        rr = r*(1 + rng.uniform(-var, var))
        p.append((cx + math.cos(a)*rr*sx, cy + math.sin(a)*rr*sy))
    return smooth_path(p)


GRAIN = '''<filter id="grain" x="0" y="0" width="100%" height="100%">
<feTurbulence type="fractalNoise" baseFrequency="0.9" numOctaves="2" seed="7" result="n"/>
<feColorMatrix type="saturate" values="0"/>
<feComponentTransfer><feFuncA type="table" tableValues="0 0.9"/></feComponentTransfer>
</filter>
<filter id="soft" x="-20%" y="-20%" width="140%" height="140%"><feGaussianBlur stdDeviation="6"/></filter>
<filter id="blur20" x="-50%" y="-50%" width="200%" height="200%"><feGaussianBlur stdDeviation="20"/></filter>
<filter id="blur60" x="-50%" y="-50%" width="200%" height="200%"><feGaussianBlur stdDeviation="60"/></filter>
<filter id="glow" x="-100%" y="-100%" width="300%" height="300%"><feGaussianBlur stdDeviation="14" result="b"/><feMerge><feMergeNode in="b"/><feMergeNode in="SourceGraphic"/></feMerge></filter>
'''


def wrap(defs, body, w=W, h=H, grain=0.07):
    g = ''
    if grain:
        g = '<rect width="%d" height="%d" filter="url(#grain)" opacity="%.2f" style="mix-blend-mode:overlay"/>' % (w, h, grain)
    return ('<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d" viewBox="0 0 %d %d">'
            '<defs>%s%s</defs>%s%s</svg>') % (w, h, w, h, GRAIN, defs, body, g)


_pw = None


def render(svg, out, w=W, h=H, quality=92):
    global _pw
    from playwright.sync_api import sync_playwright
    if _pw is None:
        p = sync_playwright().start()
        _pw = (p, p.chromium.launch())
    b = _pw[1]
    page = b.new_page(viewport={'width': w, 'height': h})
    page.set_content('<html><body style="margin:0;overflow:hidden;background:#000">%s</body></html>' % svg)
    page.wait_for_timeout(200)
    tmp = out + '.png'
    page.screenshot(path=tmp, clip={'x': 0, 'y': 0, 'width': w, 'height': h})
    page.close()
    from PIL import Image
    im = Image.open(tmp).convert('RGB')
    import os
    if out.endswith('.png'):
        os.replace(tmp, out)
    else:
        im.save(out, quality=quality, optimize=True, progressive=True)
        os.remove(tmp)
    return im


def close():
    global _pw
    if _pw:
        _pw[1].close(); _pw[0].stop(); _pw = None
