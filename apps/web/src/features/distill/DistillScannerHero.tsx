import React, { useEffect, useRef, useState } from 'react';
import './distill.css';

// 6 Hosted CDN Assets
const ASSETS = {
  dunes: 'https://d2ol7oe51mr4n9.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/0e201fcd-ffc0-46a0-bd06-faa07bf0e84b.webp',
  foundry: 'https://d2ol7oe51mr4n9.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/e80ba08d-d7b6-4b27-b0cd-8d0ad09c788c.webp',
  mineral: 'https://d2ol7oe51mr4n9.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/2bebad68-8177-446a-9550-d64074166364.webp',
  mesa: 'https://d2ol7oe51mr4n9.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/62b0d7c1-3eb7-44bd-9b19-57b975fb250f.webp',
  fire: 'https://d2ol7oe51mr4n9.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/156ff74d-ee37-4443-8f3a-11ae2095ed4c.webp',
  saltflat: 'https://d2ol7oe51mr4n9.cloudfront.net/user_38xzZboKViGWJOttwIXH07lWA1P/7b6ad638-89ec-4dec-be74-d36fbff3ea53.webp',
};

// 7-Color Fire / Hot Palette
const PALETTE = [
  '#FF4A12',
  '#FF5A1F',
  '#FF6B2C',
  '#FF8A3D',
  '#FFB06A',
  '#FFD9B5',
  '#FFF3E6',
];

interface CardMeta {
  id: string;
  original: number;
  distilled: number;
  saved: number;
  pct: string;
  url: string;
  htmlText: string;
}

const CARDS_DATA: CardMeta[] = [
  {
    id: 'home',
    original: 70411,
    distilled: 6140,
    saved: 64271,
    pct: '91% less',
    url: 'kessler.co/',
    htmlText: `<h1>Kessler — Advanced materials</h1>
<p>
  Advanced materials for a stronger
  tomorrow. Ceramic-matrix alloys made
  with desert sand and solar heat.
</p>
<h2>At a glance</h2>
<ul>
  <li>Strength: 200x steel by weight</li>
  <li>Emissions: 87% lower</li>
  <li>Plants: Nevada, Atacama, Oman</li>
  <li>Founded: 2019</li>
</ul>
<h2>Products</h2>
<ul>
  <li>K-Alloy 7 — structural</li>
  <li>Basalt fibre — tensile</li>
  <li>Ceramic matrix — thermal</li>
</ul>
<a href="/demo">Book a demo</a>`,
  },
  {
    id: 'price',
    original: 41208,
    distilled: 4960,
    saved: 36248,
    pct: '88% less',
    url: 'kessler.co/pricing',
    htmlText: `<h1>Kessler — Pricing</h1>
<p>
  Three plans for every scale of
  production. Start small, scale to
  full foundry output.
</p>
<h2>Pilot — $299 / mo</h2>
<ul>
  <li>Sample batches: 4 / mo</li>
  <li>Alloy grades: 2</li>
  <li>Support: email</li>
</ul>
<h2>Industrial — $2,400 / mo</h2>
<ul>
  <li>Tonnage: up to 40 t / mo</li>
  <li>Alloy grades: all</li>
  <li>Certification: ISO 9001</li>
  <li>Support: dedicated engineer</li>
</ul>`,
  },
  {
    id: 'case',
    original: 98432,
    distilled: 6300,
    saved: 92132,
    pct: '94% less',
    url: 'kessler.co/cases/atlas-rail',
    htmlText: `<h1>Case study: Atlas Rail</h1>
<p>
  How Atlas Rail cut track weight by
  a third in eighteen months.
</p>
<h2>The challenge</h2>
<p>
  1,200 km of desert line, rails
  warping above 50 °C every summer.
</p>
<h2>The result</h2>
<ul>
  <li>34% lighter track sections</li>
  <li>3.4x longer service life</li>
  <li>0 unplanned closures in 2025</li>
</ul>
<p>
  "We stopped losing summers."
  — Head of Infrastructure
</p>`,
  },
  {
    id: 'vision',
    original: 55870,
    distilled: 8215,
    saved: 47655,
    pct: '85% less',
    url: 'kessler.co/vision',
    htmlText: `<article>
  <h1>The Kessler vision</h1>
  <p>
    Materials that outlast the people
    who made them.
  </p>
  <p>
    We build where the sun is
    strongest and the sand is free:
    three plants, one closed loop.
  </p>
  <h2>Principles</h2>
  <ol>
    <li>Heat from the sun, not gas</li>
    <li>Sand in, scrap back in</li>
    <li>Publish every number</li>
  </ol>
</article>`,
  },
  {
    id: 'alloy',
    original: 63904,
    distilled: 5512,
    saved: 58392,
    pct: '91% less',
    url: 'kessler.co/products/k-alloy-7',
    htmlText: `<h1>K-Alloy 7</h1>
<p>The carbon-negative alloy.</p>
<h2>Specification</h2>
<ul>
  <li>Tensile strength: 1,940 MPa</li>
  <li>Density: 2.9 g/cm³</li>
  <li>Max service temp: 1,150 °C</li>
  <li>CO₂: −0.4 kg per kg made</li>
</ul>
<h2>Uses</h2>
<ul>
  <li>Rail, bridges, turbine housings</li>
</ul>
<a href="/k-alloy-7.pdf">Datasheet</a>`,
  },
  {
    id: 'journal',
    original: 47316,
    distilled: 7602,
    saved: 39714,
    pct: '84% less',
    url: 'kessler.co/journal',
    htmlText: `<h1>Kessler Journal</h1>
<p>Field notes from the flats.</p>
<article>
  <h2>Mining nothing but light</h2>
  <time>2026-03-14</time>
</article>
<article>
  <h2>Atacama, year one</h2>
  <time>2026-02-02</time>
</article>
<article>
  <h2>Why we melt sand at noon</h2>
  <time>2026-01-09</time>
</article>`,
  },
];

// 32-bit Integer Hash
function hash32(a: number, b: number, c: number): number {
  let h = (Math.imul(a, 374761393) + Math.imul(b, 668265263) + Math.imul(c, 2147483647)) | 0;
  h = Math.imul(h ^ (h >>> 13), 1274126177) | 0;
  return ((h ^ (h >>> 16)) >>> 0) / 4294967296;
}

interface Spark {
  x: number;
  y: number;
  vx: number;
  vy: number;
  size: number;
  color: string;
  life: number;
  maxLife: number;
}

interface DistillScannerHeroProps {
  onExplore?: () => void;
  onSignIn?: () => void;
  onBookDemo?: () => void;
}

export const DistillScannerHero: React.FC<DistillScannerHeroProps> = ({
  onExplore,
  onSignIn,
  onBookDemo,
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const beltRef = useRef<HTMLDivElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  const cardRefs = useRef<(HTMLDivElement | null)[]>([]);
  const codeRefs = useRef<(HTMLDivElement | null)[]>([]);

  // Interactive Counter State
  const [counterNum, setCounterNum] = useState<string>('−48,193');
  const [counterPct, setCounterPct] = useState<string>('91% less');

  // Animation Refs
  const animState = useRef({
    off: 0,
    vel: 0,
    hover: 0,
    isHovered: false,
    isDragging: false,
    dragStartX: 0,
    dragStartOff: 0,
    lastMoveX: 0,
    lastMoveT: 0,
    targetNum: 48193,
    currentNum: 48193,
    lastCrossingIndex: 0,
    sparks: [] as Spark[],
  });

  useEffect(() => {
    const beltEl = beltRef.current;
    const canvas = canvasRef.current;
    if (!beltEl || !canvas) return;

    const ctx = canvas.getContext('2d');
    if (!ctx) return;

    let rafId: number;
    let lastT = performance.now();

    // Check reduced motion preference
    const mediaReducedMotion = window.matchMedia('(prefers-reduced-motion: reduce)');
    let prefersReducedMotion = mediaReducedMotion.matches;
    const handleMotionChange = (e: MediaQueryListEvent) => {
      prefersReducedMotion = e.matches;
    };
    mediaReducedMotion.addEventListener('change', handleMotionChange);

    // Initial Geometry Setup
    const updateGeometry = () => {
      const W = window.innerWidth;
      const dpr = Math.min(window.devicePixelRatio || 1, 2);
      const beltRect = beltEl.getBoundingClientRect();

      canvas.width = beltRect.width * dpr;
      canvas.height = beltRect.height * dpr;
      canvas.style.width = `${beltRect.width}px`;
      canvas.style.height = `${beltRect.height}px`;

      // Set initial offset: W/2 - 0.62 * cw so first card is ~12% past line
      const firstCard = cardRefs.current[0];
      const cw = firstCard ? firstCard.getBoundingClientRect().width : W * 0.46;
      if (animState.current.off === 0) {
        animState.current.off = W / 2 - 0.62 * cw;
      }
    };

    updateGeometry();
    window.addEventListener('resize', updateGeometry);

    // Belt Pointer / Drag Listeners
    const onPointerDown = (e: PointerEvent) => {
      animState.current.isDragging = true;
      animState.current.dragStartX = e.clientX;
      animState.current.dragStartOff = animState.current.off;
      animState.current.lastMoveX = e.clientX;
      animState.current.lastMoveT = performance.now();
      beltEl.classList.add('is-dragging');
      beltEl.setPointerCapture(e.pointerId);
    };

    const onPointerMove = (e: PointerEvent) => {
      if (!animState.current.isDragging) return;
      const now = performance.now();
      const dt = Math.max(now - animState.current.lastMoveT, 8);
      const dx = e.clientX - animState.current.lastMoveX;

      animState.current.off = animState.current.dragStartOff + (e.clientX - animState.current.dragStartX);

      // Estimate velocity with smoothing
      const instantVel = (dx / dt) * 1000;
      animState.current.vel = instantVel * 0.8 + animState.current.vel * 0.2;

      animState.current.lastMoveX = e.clientX;
      animState.current.lastMoveT = now;
    };

    const onPointerUp = (e: PointerEvent) => {
      if (animState.current.isDragging) {
        animState.current.isDragging = false;
        beltEl.classList.remove('is-dragging');
        try {
          beltEl.releasePointerCapture(e.pointerId);
        } catch {}
      }
    };

    const onPointerEnter = () => {
      animState.current.isHovered = true;
    };
    const onPointerLeave = () => {
      animState.current.isHovered = false;
    };

    const onWheel = (e: WheelEvent) => {
      if (Math.abs(e.deltaX) > Math.abs(e.deltaY)) {
        e.preventDefault();
        animState.current.vel += -e.deltaX * 6;
      }
    };

    beltEl.addEventListener('pointerdown', onPointerDown);
    window.addEventListener('pointermove', onPointerMove);
    window.addEventListener('pointerup', onPointerUp);
    beltEl.addEventListener('pointerenter', onPointerEnter);
    beltEl.addEventListener('pointerleave', onPointerLeave);
    beltEl.addEventListener('wheel', onWheel, { passive: false });

    // Main Animation Loop
    const loop = (now: number) => {
      const dt = Math.min((now - lastT) / 1000, 0.1);
      lastT = now;

      const state = animState.current;
      const W = window.innerWidth;
      const Lx = W / 2;
      const dpr = Math.min(window.devicePixelRatio || 1, 2);

      const firstCard = cardRefs.current[0];
      const cw = firstCard ? firstCard.getBoundingClientRect().width : W * 0.46;
      const ch = firstCard ? firstCard.getBoundingClientRect().height : window.innerHeight * 0.52;
      const gap = W <= 767 ? W * 0.07 : W * 0.0475;
      const span = cw + gap;
      const period = span * 6;

      // Speed calculation
      const autoSpeed = prefersReducedMotion ? 0 : -0.075 * W;

      // Hover easing
      const targetHover = state.isHovered ? 1 : 0;
      state.hover += (targetHover - state.hover) * Math.min(dt * 4, 1);

      if (!state.isDragging) {
        const targetVel = autoSpeed * (1 - 0.65 * state.hover);
        // Ease velocity towards target at 1 - e^(-2.2 * dt)
        const easeFactor = 1 - Math.exp(-2.2 * dt);
        state.vel += (targetVel - state.vel) * easeFactor;
        state.off += state.vel * dt;
      }

      // Modulo wrap helper
      const mod = (n: number, m: number) => ((n % m) + m) % m;

      // Clear Canvas
      ctx.save();
      ctx.scale(dpr, dpr);
      ctx.clearRect(0, 0, canvas.width / dpr, canvas.height / dpr);

      let activeCrossingFound = false;

      // Position Cards & Twins
      for (let i = 0; i < 6; i++) {
        const cardEl = cardRefs.current[i];
        const codeEl = codeRefs.current[i];
        if (!cardEl || !codeEl) continue;

        const x = mod(state.off + i * span + span, period) - span;
        const cut = Math.min(Math.max(Lx - x, 0), cw);

        // Position Card
        cardEl.style.transform = `translate3d(${x}px, 0, 0)`;
        if (cut > 0) {
          cardEl.style.clipPath = `inset(0 0 0 ${cut}px round 0 1.6vw 1.6vw 0)`;
          cardEl.style.maskImage = `linear-gradient(90deg, transparent ${cut}px, #000 ${cut + 26}px)`;
          (cardEl.style as any).webkitMaskImage = `linear-gradient(90deg, transparent ${cut}px, #000 ${cut + 26}px)`;
        } else {
          cardEl.style.clipPath = 'none';
          cardEl.style.maskImage = 'none';
          (cardEl.style as any).webkitMaskImage = 'none';
        }

        // Position Twin Code Block
        codeEl.style.transform = `translate3d(${x}px, 0, 0)`;
        if (cut > 0) {
          codeEl.style.visibility = 'visible';
          codeEl.style.clipPath = `inset(0 ${cw - cut}px 0 0)`;
        } else {
          codeEl.style.visibility = 'hidden';
        }

        // Check Counter for crossing card
        if (cut > 0 && cut < cw) {
          activeCrossingFound = true;
          const meta = CARDS_DATA[i];
          const progress = cut / cw;
          const target = Math.round((meta.original - meta.distilled) * progress);

          if (state.lastCrossingIndex !== i) {
            state.lastCrossingIndex = i;
            state.currentNum = target;
            setCounterPct(meta.pct);
          }
          state.targetNum = target;
        }

        // Canvas Dissolve Effect
        if (cut > 0 && x + cw > Lx) {
          const cardTop = 70; // top offset in belt
          const hazeW = Math.min(150, x + cw - Lx);

          // 1. Haze
          if (hazeW > 0) {
            const hazeGrad = ctx.createLinearGradient(Lx, 0, Lx + hazeW, 0);
            hazeGrad.addColorStop(0, 'rgba(255, 96, 40, 0.42)');
            hazeGrad.addColorStop(1, 'transparent');
            ctx.fillStyle = hazeGrad;
            ctx.fillRect(Lx, cardTop, hazeW, ch);
          }

          // 2. Dithered Particles Fixed to Card Surface
          const R = 0.3 * cw;
          const startCol = Math.floor(cut / 3);
          const endCol = Math.floor((cut + R) / 3);
          const numRows = Math.floor(ch / 3);
          const tSec = now / 1000;

          for (let gx = startCol; gx <= endCol; gx++) {
            const sx = x + gx * 3;
            const d = sx - Lx;
            if (d < 0 || d > R) continue;

            const p = Math.pow(1 - d / R, 1.7);

            for (let row = -5; row <= numRows + 5; row++) {
              // Twinkle check
              if (!prefersReducedMotion) {
                const tw = hash32(gx, row, Math.floor(tSec * 9) + 97);
                if (tw < 0.22) continue;
              }

              const h = hash32(gx, row, i);
              const isInside = row >= 0 && row < numRows;
              const threshold = isInside
                ? p * 0.92
                : p * 0.16 * (1 - Math.abs(row < 0 ? row : row - numRows) / 6);

              if (h < threshold) {
                const h2 = hash32(row, gx, i + 31);
                // Color selection
                let color: string;
                if (d <= 10) {
                  // White-hot edge
                  color = PALETTE[4 + Math.floor(h2 * 3)] || PALETTE[6];
                } else {
                  color = PALETTE[Math.floor(h2 * h2 * 7)] || PALETTE[0];
                }

                // Size & Jitter
                const size = 1.1 + h2 * 1.7;
                const jitterX = prefersReducedMotion ? 0 : (h - 0.5) * 1.6;
                const jitterY = prefersReducedMotion ? 0 : (h2 - 0.5) * 1.6;

                const py = cardTop + row * 3 + jitterY;
                const px = sx + jitterX;

                ctx.fillStyle = color;
                ctx.fillRect(px, py, size, size);
              }
            }
          }
        }
      }

      // 3. Free Sparks
      if (!prefersReducedMotion) {
        // Spawn sparks
        if (Math.random() < 0.7) {
          const spawnY = 70 - 10 + Math.random() * (ch + 20);
          const spawnX = Lx + Math.random() * 30;
          const maxLife = 0.3 + Math.random() * 0.8;
          state.sparks.push({
            x: spawnX,
            y: spawnY,
            vx: -30 + Math.random() * 70 + state.vel * 0.6,
            vy: -25 + Math.random() * 50,
            size: 0.8 + Math.random() * 1.8,
            color: PALETTE[Math.floor(Math.random() * PALETTE.length)],
            life: maxLife,
            maxLife: maxLife,
          });
        }

        // Update and render sparks
        for (let sIdx = state.sparks.length - 1; sIdx >= 0; sIdx--) {
          const sp = state.sparks[sIdx];
          sp.life -= dt;
          if (sp.life <= 0) {
            state.sparks.splice(sIdx, 1);
            continue;
          }
          sp.x += sp.vx * dt;
          sp.y += sp.vy * dt;

          const alpha = Math.min(1, sp.life * 2.2);
          ctx.save();
          ctx.globalAlpha = alpha;
          ctx.fillStyle = sp.color;
          ctx.fillRect(sp.x, sp.y, sp.size, sp.size);
          ctx.restore();
        }
      }

      ctx.restore();

      // Smooth Counter Easing
      if (activeCrossingFound) {
        state.currentNum += (state.targetNum - state.currentNum) * Math.min(dt * 14, 1);
        const formatted = '−' + Math.round(state.currentNum).toLocaleString('en-US');
        setCounterNum(formatted);
      }

      rafId = requestAnimationFrame(loop);
    };

    rafId = requestAnimationFrame(loop);

    return () => {
      cancelAnimationFrame(rafId);
      window.removeEventListener('resize', updateGeometry);
      beltEl.removeEventListener('pointerdown', onPointerDown);
      window.removeEventListener('pointermove', onPointerMove);
      window.removeEventListener('pointerup', onPointerUp);
      beltEl.removeEventListener('pointerenter', onPointerEnter);
      beltEl.removeEventListener('pointerleave', onPointerLeave);
      beltEl.removeEventListener('wheel', onWheel);
      mediaReducedMotion.removeEventListener('change', handleMotionChange);
    };
  }, []);

  return (
    <div className="distill-stage" ref={containerRef}>
      {/* 56px Nav Bar */}
      <header className="distill-nav">
        <a className="distill-logo" href="#distill">
          <svg className="distill-logo-icon" viewBox="0 0 19 23">
            <path d="M9.5 0C9.5 0 0 11.5 0 16.5C0 20.0899 4.25329 23 9.5 23C14.7467 23 19 20.0899 19 16.5C19 11.5 9.5 0 9.5 0Z" />
            <path
              d="M7 19.5C5.5 18.5 4.5 16.8 4.5 15"
              fill="none"
              stroke="#EEF0F2"
              strokeWidth="1.7"
              strokeLinecap="round"
            />
          </svg>
          Distill
        </a>

        <nav className="distill-links">
          <a className="distill-link" href="#product">
            Product
            <svg width="8" height="5" viewBox="0 0 8 5" fill="none">
              <path d="M1 1L4 4L7 1" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
          </a>
          <a className="distill-link" href="#use-cases">
            Use cases
            <svg width="8" height="5" viewBox="0 0 8 5" fill="none">
              <path d="M1 1L4 4L7 1" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
          </a>
          <a className="distill-link" href="#resources">
            Resources
            <svg width="8" height="5" viewBox="0 0 8 5" fill="none">
              <path d="M1 1L4 4L7 1" stroke="currentColor" strokeWidth="1.2" strokeLinecap="round" strokeLinejoin="round"/>
            </svg>
          </a>
          <a className="distill-link" href="#pricing">Pricing</a>
        </nav>

        <div className="distill-nav-right">
          <button className="distill-btn-signin" onClick={onSignIn}>Sign in</button>
          <button className="distill-btn-trial" onClick={onExplore}>Start a free trial</button>
          <button className="distill-btn-demo" onClick={onBookDemo}>Book a demo</button>
          <button className="distill-btn-app" onClick={onExplore} title="Open PayBridge 16-Screen Enterprise Suite">
            PayBridge Suite ↗
          </button>
        </div>
      </header>

      {/* Hero Copy */}
      <section className="distill-hero">
        <div className="distill-copy-left">
          <div className="distill-eyebrow">AI READABILITY</div>
          <h1 className="distill-headline">
            Built for people.
            <em>Read by machines.</em>
          </h1>
        </div>

        <div className="distill-copy-right">
          <p className="distill-paragraph">
            Distill detects AI agents at the edge and serves them a clean, token‑light version of every page, without changing a pixel of what people see.
          </p>
          <button className="distill-explore-btn" onClick={onExplore}>
            Explore Distill
          </button>
        </div>
      </section>

      {/* Floating Counter Badge on Top of Line */}
      <div className="distill-badge">
        <span className="distill-badge-label">Token reduction</span>
        <div className="distill-badge-pill">
          <span className="distill-badge-num">{counterNum}</span>
          <span className="distill-badge-pct">{counterPct}</span>
        </div>
      </div>

      {/* Glowing Scanner Line */}
      <div className="distill-scan">
        <i className="distill-scan-core" />
        <b className="distill-scan-halo" />
      </div>

      {/* The Conveyor Belt */}
      <div className="distill-belt" ref={beltRef}>
        {/* Layer 1: Code Layer */}
        <div className="distill-codes">
          {CARDS_DATA.map((card, idx) => (
            <div
              key={`code-${card.id}`}
              className="k-code"
              ref={(el) => { codeRefs.current[idx] = el; }}
            >
              <div className="k-code-html">{card.htmlText}</div>
              <div className="k-code-meta">
                {`<!-- distilled: ${card.distilled.toLocaleString('en-US')} tokens -->`}
                <br />
                {`<!-- original: ${card.original.toLocaleString('en-US')} tokens -->`}
                <br />
                {`<!-- ${card.url} -->`}
              </div>
            </div>
          ))}
        </div>

        {/* Layer 2: Website Cards */}
        <div className="distill-cards">
          {/* Card 1: Home (Dunes) */}
          <div
            className="k-card"
            ref={(el) => { cardRefs.current[0] = el; }}
            style={{
              backgroundImage: `linear-gradient(180deg, rgba(20,10,4,0.25) 0%, transparent 30%, rgba(20,10,4,0.55) 100%), url(${ASSETS.dunes})`,
              backgroundSize: 'cover',
              backgroundPosition: 'center',
            }}
          >
            <div className="k-nav">
              <div className="k-nav-logo">
                <svg viewBox="0 0 54 42" fill="none">
                  <path d="M2 38L18 6L32 38H2Z" stroke="#fff" strokeWidth="3.2" strokeLinejoin="round" />
                  <path d="M22 38L38 12L52 38H22Z" stroke="#fff" strokeWidth="3.2" strokeLinejoin="round" />
                </svg>
                <span>Kessler</span>
              </div>
              <div className="k-nav-menu">
                <span className="k-nav-pill">PRODUCTS</span>
                <span>IMPACT</span>
                <span>ABOUT</span>
              </div>
            </div>

            <div style={{ position: 'absolute', bottom: '4.5cqw', left: '4.4cqw', display: 'flex', gap: '4cqw', alignItems: 'flex-end' }}>
              <div>
                <div style={{ fontSize: '9cqw', lineHeight: 0.9, letterSpacing: '-0.04em', fontFamily: 'var(--font-sans)', fontWeight: 300 }}>200x</div>
                <div style={{ fontSize: '1.9cqw', opacity: 0.85, marginTop: '0.4cqw' }}>stronger</div>
              </div>
              <div style={{ width: '1px', height: '9cqw', background: 'rgba(255,255,255,0.35)' }} />
              <div>
                <div style={{ fontSize: '9cqw', lineHeight: 0.9, letterSpacing: '-0.04em', fontFamily: 'var(--font-sans)', fontWeight: 300 }}>87%</div>
                <div style={{ fontSize: '1.9cqw', opacity: 0.85, marginTop: '0.4cqw' }}>lower emissions</div>
              </div>
            </div>

            <div style={{ position: 'absolute', bottom: '4.5cqw', right: '4.4cqw', width: '24cqw', fontSize: '2.2cqw', lineHeight: 1.25, textAlign: 'right', opacity: 0.9 }}>
              Advanced materials for a stronger tomorrow.
            </div>
          </div>

          {/* Card 2: Price (Radial Gradient) */}
          <div
            className="k-card"
            ref={(el) => { cardRefs.current[1] = el; }}
            style={{
              background: 'radial-gradient(120% 90% at 20% 110%, #FF6A1F 0%, #FF9A4A 35%, #F7C79A 62%, #EFD9C4 100%)',
              color: '#2A1406',
            }}
          >
            <div className="k-nav">
              <div className="k-nav-logo" style={{ color: '#2A1406' }}>
                <svg viewBox="0 0 54 42" fill="none">
                  <path d="M2 38L18 6L32 38H2Z" stroke="#2A1406" strokeWidth="3.2" strokeLinejoin="round" />
                  <path d="M22 38L38 12L52 38H22Z" stroke="#2A1406" strokeWidth="3.2" strokeLinejoin="round" />
                </svg>
                <span>Kessler</span>
              </div>
              <div className="k-nav-menu" style={{ color: 'rgba(42,20,6,0.85)' }}>
                <span>PRODUCTS</span>
                <span className="k-nav-pill" style={{ borderColor: 'rgba(42,20,6,0.5)', color: '#2A1406' }}>PRICING</span>
                <span>CONTACT</span>
              </div>
            </div>

            <div style={{ position: 'absolute', top: '38%', left: '4.4cqw', fontSize: '13cqw', lineHeight: 0.9, letterSpacing: '-0.05em', fontWeight: 300, fontFamily: 'var(--font-sans)' }}>
              <sup style={{ fontSize: '0.4em', top: '-0.4em' }}>$</sup>2,400
            </div>
            <div style={{ position: 'absolute', top: '62%', left: '4.4cqw', fontFamily: 'var(--font-mono)', fontSize: '1.6cqw', letterSpacing: '0.08em', opacity: 0.75 }}>
              PER MONTH · INDUSTRIAL
            </div>

            <div style={{ position: 'absolute', top: '36%', right: '4.4cqw', width: '36cqw', display: 'flex', flexDirection: 'column' }}>
              <div style={{ borderTop: '1px solid rgba(42,20,6,0.25)', padding: '1.2cqw 0', display: 'flex', justifyContent: 'space-between', fontSize: '1.9cqw' }}>
                <span>Pilot</span><span style={{ fontWeight: 500 }}>$299</span>
              </div>
              <div style={{ borderTop: '1px solid rgba(42,20,6,0.25)', padding: '1.2cqw 0', display: 'flex', justifyContent: 'space-between', fontSize: '1.9cqw' }}>
                <span>Industrial</span><span style={{ fontWeight: 500 }}>$2,400</span>
              </div>
              <div style={{ borderTop: '1px solid rgba(42,20,6,0.25)', padding: '1.2cqw 0', display: 'flex', justifyContent: 'space-between', fontSize: '1.9cqw' }}>
                <span>Foundry</span><span style={{ fontWeight: 500 }}>Custom</span>
              </div>
            </div>

            <button
              style={{
                position: 'absolute',
                bottom: '4.5cqw',
                left: '4.4cqw',
                padding: '1.2cqw 3cqw',
                borderRadius: '9999px',
                background: '#2A1406',
                color: '#F7E6D6',
                border: 'none',
                fontSize: '1.9cqw',
                fontWeight: 500,
                cursor: 'pointer',
              }}
            >
              Start a pilot
            </button>
          </div>

          {/* Card 3: Case Study (Mineral) */}
          <div
            className="k-card"
            ref={(el) => { cardRefs.current[2] = el; }}
            style={{
              backgroundImage: `linear-gradient(180deg, rgba(25,14,8,0.15) 0%, rgba(25,14,8,0.55) 100%), url(${ASSETS.mineral})`,
              backgroundSize: '112%',
              backgroundPosition: 'center',
              filter: 'saturate(1.1)',
            }}
          >
            <div className="k-nav">
              <div className="k-nav-logo">
                <svg viewBox="0 0 54 42" fill="none">
                  <path d="M2 38L18 6L32 38H2Z" stroke="#fff" strokeWidth="3.2" strokeLinejoin="round" />
                  <path d="M22 38L38 12L52 38H22Z" stroke="#fff" strokeWidth="3.2" strokeLinejoin="round" />
                </svg>
                <span>Kessler</span>
              </div>
              <div className="k-nav-menu">
                <span>PRODUCTS</span>
                <span className="k-nav-pill">IMPACT</span>
                <span>ABOUT</span>
              </div>
            </div>

            <div style={{ position: 'absolute', top: '17%', left: '4.4cqw', fontFamily: 'var(--font-mono)', fontSize: '1.4cqw', letterSpacing: '0.08em', opacity: 0.85 }}>
              CASE STUDY · ATLAS RAIL
            </div>
            <div style={{ position: 'absolute', top: '24%', left: '4.4cqw', width: '56cqw', fontFamily: 'var(--font-serif)', fontSize: '4.6cqw', lineHeight: 1.15 }}>
              How Atlas Rail cut track weight by a third in eighteen months
            </div>

            <div style={{ position: 'absolute', bottom: '4.5cqw', left: '4.4cqw', right: '4.4cqw', display: 'grid', gridTemplateColumns: 'repeat(3, 1fr)', gap: '2cqw' }}>
              <div style={{ borderTop: '1px solid rgba(255,255,255,0.35)', paddingTop: '1cqw' }}>
                <div style={{ fontSize: '6cqw', fontWeight: 300, lineHeight: 1 }}>34%</div>
                <div style={{ fontSize: '1.6cqw', opacity: 0.85 }}>lighter track</div>
              </div>
              <div style={{ borderTop: '1px solid rgba(255,255,255,0.35)', paddingTop: '1cqw' }}>
                <div style={{ fontSize: '6cqw', fontWeight: 300, lineHeight: 1 }}>3.4x</div>
                <div style={{ fontSize: '1.6cqw', opacity: 0.85 }}>service life</div>
              </div>
              <div style={{ borderTop: '1px solid rgba(255,255,255,0.35)', paddingTop: '1cqw' }}>
                <div style={{ fontSize: '6cqw', fontWeight: 300, lineHeight: 1 }}>0</div>
                <div style={{ fontSize: '1.6cqw', opacity: 0.85 }}>unplanned closures</div>
              </div>
            </div>
          </div>

          {/* Card 4: Vision (Foundry & Mesa Tiles) */}
          <div
            className="k-card"
            ref={(el) => { cardRefs.current[3] = el; }}
            style={{
              background: '#1B1311',
              color: '#fff',
            }}
          >
            <div className="k-nav">
              <div className="k-nav-logo">
                <svg viewBox="0 0 54 42" fill="none">
                  <path d="M2 38L18 6L32 38H2Z" stroke="#fff" strokeWidth="3.2" strokeLinejoin="round" />
                  <path d="M22 38L38 12L52 38H22Z" stroke="#fff" strokeWidth="3.2" strokeLinejoin="round" />
                </svg>
                <span>Kessler</span>
              </div>
              <div className="k-nav-menu">
                <span>PRODUCTS</span>
                <span className="k-nav-pill">VISION</span>
                <span>ABOUT</span>
              </div>
            </div>

            <div style={{ position: 'absolute', top: '21%', left: '4.4cqw', right: '4.4cqw', height: '46%', display: 'grid', gridTemplateColumns: '1fr 1fr 1.2fr', gap: '2cqw' }}>
              <div style={{ borderRadius: '1.4cqw', backgroundImage: `url(${ASSETS.foundry})`, backgroundSize: 'cover', backgroundPosition: 'center' }} />
              <div style={{ borderRadius: '1.4cqw', backgroundImage: `url(${ASSETS.mesa})`, backgroundSize: 'cover', backgroundPosition: 'center' }} />
              <div style={{ display: 'flex', alignItems: 'center', padding: '0 1cqw', fontFamily: 'var(--font-serif)', fontSize: '3.4cqw', lineHeight: 1.25 }}>
                <span>The Kessler vision: <em style={{ opacity: 0.75 }}>materials that outlast the people who made them.</em></span>
              </div>
            </div>

            <div style={{ position: 'absolute', bottom: '4.5cqw', left: '4.4cqw', right: '4.4cqw', display: 'flex', justifyContent: 'space-between', fontSize: '1.6cqw', opacity: 0.7 }}>
              <span>How Kessler pours alloy / at the edge of the desert</span>
              <span>Manifesto · 6 min</span>
            </div>
          </div>

          {/* Card 5: Alloy (Fire) */}
          <div
            className="k-card"
            ref={(el) => { cardRefs.current[4] = el; }}
            style={{
              backgroundImage: `linear-gradient(90deg, rgba(20,12,8,0.3) 0%, rgba(20,12,8,0.7) 100%), url(${ASSETS.fire})`,
              backgroundSize: '106%',
              backgroundPosition: 'center',
              backgroundColor: '#140C08',
            }}
          >
            <div className="k-nav">
              <div className="k-nav-logo">
                <svg viewBox="0 0 54 42" fill="none">
                  <path d="M2 38L18 6L32 38H2Z" stroke="#fff" strokeWidth="3.2" strokeLinejoin="round" />
                  <path d="M22 38L38 12L52 38H22Z" stroke="#fff" strokeWidth="3.2" strokeLinejoin="round" />
                </svg>
                <span>Kessler</span>
              </div>
              <div className="k-nav-menu">
                <span className="k-nav-pill">PRODUCTS</span>
                <span>IMPACT</span>
                <span>ABOUT</span>
              </div>
            </div>

            <div style={{ position: 'absolute', bottom: '4.5cqw', left: '4.4cqw', fontFamily: 'var(--font-serif)', fontSize: '7.4cqw', lineHeight: 0.98 }}>
              <div>The carbon‑negative</div>
              <em>alloy.</em>
            </div>

            <div style={{ position: 'absolute', bottom: '4.5cqw', right: '4.4cqw', width: '30cqw', display: 'flex', flexDirection: 'column' }}>
              <div style={{ borderTop: '1px solid rgba(255,255,255,0.3)', padding: '1cqw 0', display: 'flex', justifyContent: 'space-between', fontSize: '1.8cqw' }}>
                <span style={{ opacity: 0.75 }}>Product</span><span>K-Alloy 7</span>
              </div>
              <div style={{ borderTop: '1px solid rgba(255,255,255,0.3)', padding: '1cqw 0', display: 'flex', justifyContent: 'space-between', fontSize: '1.8cqw' }}>
                <span style={{ opacity: 0.75 }}>Tensile</span><span>1,940 MPa</span>
              </div>
              <div style={{ borderTop: '1px solid rgba(255,255,255,0.3)', padding: '1cqw 0', display: 'flex', justifyContent: 'space-between', fontSize: '1.8cqw' }}>
                <span style={{ opacity: 0.75 }}>Datasheet</span><span>PDF ↗</span>
              </div>
            </div>
          </div>

          {/* Card 6: Journal (Light #E9ECF1) */}
          <div
            className="k-card"
            ref={(el) => { cardRefs.current[5] = el; }}
            style={{
              background: '#E9ECF1',
              color: '#121A2A',
            }}
          >
            <div className="k-nav">
              <div className="k-nav-logo" style={{ color: '#121A2A' }}>
                <svg viewBox="0 0 54 42" fill="none">
                  <path d="M2 38L18 6L32 38H2Z" stroke="#121A2A" strokeWidth="3.2" strokeLinejoin="round" />
                  <path d="M22 38L38 12L52 38H22Z" stroke="#121A2A" strokeWidth="3.2" strokeLinejoin="round" />
                </svg>
                <span>Kessler</span>
              </div>
              <div className="k-nav-menu" style={{ color: 'rgba(18,26,42,0.7)' }}>
                <span>PRODUCTS</span>
                <span className="k-nav-pill" style={{ borderColor: 'rgba(18,26,42,0.4)', color: '#121A2A' }}>JOURNAL</span>
                <span>ABOUT</span>
              </div>
            </div>

            <div style={{ position: 'absolute', top: '23%', left: '4.4cqw', fontFamily: 'var(--font-serif)', fontSize: '5.6cqw', lineHeight: 1.05 }}>
              <div>Journal</div>
              <em style={{ color: '#46618E', fontSize: '0.85em' }}>field notes from the flats</em>
            </div>

            <div style={{ position: 'absolute', bottom: '4.5cqw', left: '4.4cqw', right: '4.4cqw', display: 'grid', gridTemplateColumns: '1.25fr 1fr 1fr', gap: '2.5cqw' }}>
              <div>
                <div style={{ aspectRatio: '4/3', borderRadius: '1.2cqw', backgroundImage: `url(${ASSETS.saltflat})`, backgroundSize: 'cover', backgroundPosition: 'center', marginBottom: '1cqw' }} />
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1.3cqw', color: '#6A7891' }}>FIELD NOTES</div>
                <div style={{ fontSize: '1.6cqw', fontWeight: 500, lineHeight: 1.25 }}>Mining nothing but light on the salt flats</div>
              </div>

              <div>
                <div style={{ aspectRatio: '4/3', borderRadius: '1.2cqw', backgroundImage: `url(${ASSETS.mesa})`, backgroundSize: 'cover', backgroundPosition: 'center', marginBottom: '1cqw' }} />
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1.3cqw', color: '#6A7891' }}>PLANTS</div>
                <div style={{ fontSize: '1.6cqw', fontWeight: 500, lineHeight: 1.25 }}>Atacama, year one</div>
              </div>

              <div>
                <div style={{ aspectRatio: '4/3', borderRadius: '1.2cqw', backgroundImage: `url(${ASSETS.dunes})`, backgroundSize: 'cover', backgroundPosition: 'center', marginBottom: '1cqw' }} />
                <div style={{ fontFamily: 'var(--font-mono)', fontSize: '1.3cqw', color: '#6A7891' }}>PROCESS</div>
                <div style={{ fontSize: '1.6cqw', fontWeight: 500, lineHeight: 1.25 }}>Why we melt sand at noon</div>
              </div>
            </div>
          </div>
        </div>

        {/* Layer 3: 2D Particle Dissolve Canvas */}
        <canvas className="distill-canvas" ref={canvasRef} />
      </div>
    </div>
  );
};
