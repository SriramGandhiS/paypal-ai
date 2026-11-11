import React, { useEffect, useRef } from 'react';
import Hls from 'hls.js';

const HLS_SRC = 'https://stream.mux.com/bnYL6x5cAX6WiJv2pOKpITehZd3NVdXpj3ylJFpX5Lk.m3u8';

const CROSS_ICON = 'https://cdn.prod.website-files.com/6720dd1ab6df0da205830ab1/686cc0f520a992816d8b15dc_bullet-list-cross.svg';
const CHECK_ICON = 'https://cdn.prod.website-files.com/6720dd1ab6df0da205830ab1/686cc068490683bbb3377d04_bullet-list.svg';

const negatives = [
  'Reactive firefighting when foundational issues surface too late',
  'Bloated coordination overhead drains bandwidth from core teams',
  "Constant re-verification because source data can't be trusted",
  'Fragmented vendor relations produce mismatched deliverables',
  'Scattered specs and decisions buried across siloed systems',
];

const positives = [
  'Layered dependency maps eliminate costly surprises at every phase',
  'Streamlined team handoffs deliver production-ready outcomes fast',
  'Live validation loops keep requirements locked across all stages',
  'Unified vendor management through a single accountable contact',
  'Centralized context and clear records accelerate every decision',
];

function HlsVideo() {
  const videoRef = useRef<HTMLVideoElement>(null);

  useEffect(() => {
    const video = videoRef.current;
    if (!video) return;

    if (Hls.isSupported()) {
      const hls = new Hls({
        startLevel: -1,
        capLevelToPlayerSize: false,
        maxMaxBufferLength: 60,
        enableWorker: true,
      });
      hls.loadSource(HLS_SRC);
      hls.attachMedia(video);
      hls.on(Hls.Events.MANIFEST_PARSED, () => {
        hls.currentLevel = hls.levels.length - 1;
        video.play().catch(() => {});
      });
      return () => hls.destroy();
    } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
      video.src = HLS_SRC;
      video.play().catch(() => {});
    }
  }, []);

  return (
    <video
      ref={videoRef}
      autoPlay
      loop
      muted
      playsInline
      style={{
        width: '160%',
        height: '160%',
        objectFit: 'cover',
        position: 'absolute',
        top: '50%',
        left: '50%',
        transform: 'translate(-50%, -50%)',
      }}
    />
  );
}

export const FreedomSection: React.FC = () => {
  return (
    <section
      className="w-full flex flex-col items-center"
      style={{
        backgroundColor: '#ffffff',
        padding: 'clamp(48px, 6vw, 80px) clamp(16px, 3vw, 40px)',
        gap: '36px',
      }}
    >
      {/* Block 1 — Header */}
      <div className="flex flex-col items-center gap-9 text-center">
        {/* Badge Pill */}
        <div
          className="flex items-center gap-2 text-lg font-medium rounded-full"
          style={{
            backgroundColor: 'rgb(249, 249, 249)',
            padding: '0.9vw 1.25vw',
            color: 'rgb(26, 11, 84)',
          }}
        >
          <svg
            width="19"
            height="18"
            style={{ flexShrink: 0 }}
            viewBox="0 0 17 16"
            fill="none"
            xmlns="http://www.w3.org/2000/svg"
          >
            <g clipPath="url(#freedom-clip)">
              <path
                fillRule="evenodd"
                clipRule="evenodd"
                d="M8.50037 3.66955C7.53221 2.82462 6.41758 2.275 5.333 2.07887C4.11096 1.85888 2.84987 2.0826 1.96658 2.95885C1.10056 3.81944 0.866218 5.04172 1.06751 6.23193C1.24778 7.29835 1.7803 8.39907 2.60501 9.35959C2.41536 10.1071 2.46371 10.8946 2.7434 11.6137C3.02308 12.3327 3.52035 12.9481 4.16678 13.375C4.81321 13.802 5.57702 14.0195 6.35308 13.9976C7.12915 13.9758 7.87933 13.7157 8.50037 13.2531C9.12146 13.7161 9.87183 13.9765 10.6482 13.9985C11.4245 14.0205 12.1886 13.8029 12.8352 13.3758C13.4819 12.9487 13.9792 12.3331 14.2588 11.6137C14.5384 10.8943 14.5865 10.1065 14.3965 9.35884C15.2204 8.39832 15.753 7.29835 15.9325 6.23119C16.1338 5.04098 15.8994 3.81944 15.0334 2.9596C14.1501 2.0826 12.889 1.85888 11.667 2.07962C10.5824 2.275 9.46854 2.82537 8.50037 3.66955Z"
                fill="rgb(200, 111, 255)"
              />
            </g>
            <defs>
              <clipPath id="freedom-clip">
                <rect width="16" height="16" fill="white" transform="translate(0.5)" />
              </clipPath>
            </defs>
          </svg>
          <span>Control</span>
        </div>

        {/* Headline */}
        <h2
          className="font-medium"
          style={{
            fontSize: 'clamp(32px, 4vw, 56px)',
            color: 'rgb(26, 11, 84)',
            lineHeight: 1.15,
            margin: 0,
          }}
        >
          Stop absorbing the chaos.
          <br />
          <span
            style={{
              backgroundImage:
                'linear-gradient(90deg, rgb(43,167,255), rgb(202,69,255) 50%, rgb(254,136,27))',
              WebkitBackgroundClip: 'text',
              backgroundClip: 'text',
              WebkitTextFillColor: 'transparent',
              color: 'transparent',
              paddingBottom: '0.3vw',
              display: 'inline-block',
            }}
          >
            Run with confidence.
          </span>
        </h2>
      </div>

      {/* Block 2 — Three-column grid */}
      <div
        className="w-full flex flex-col lg:grid"
        style={{
          gridTemplateColumns: '26vw 1fr 26vw',
          columnGap: '36px',
          rowGap: '24px',
          alignItems: 'start',
          padding: '0 clamp(0px, 2.92vw, 40px)',
          gap: '24px',
        }}
      >
        {/* Left column — Negatives */}
        <div
          className="flex flex-col"
          style={{
            gap: '12px',
            fontSize: 'clamp(13px, 1.15vw, 17px)',
            color: 'rgb(131, 121, 158)',
          }}
        >
          {negatives.map((text, i) => (
            <div
              key={i}
              className="flex flex-col"
              style={{
                gap: '12px',
                padding: 'clamp(12px, 0.97vw, 16px) clamp(14px, 1.25vw, 20px)',
                borderRadius: '18px',
                backgroundColor: 'rgb(255, 255, 255)',
                boxShadow: '0 3px 9.1px #3f4a7e0d, 0 1px 29px #3f4a7e1a',
              }}
            >
              <img
                src={CROSS_ICON}
                alt=""
                aria-hidden
                style={{ width: 'clamp(16px, 1.25vw, 20px)', flexShrink: 0 }}
              />
              <div>{text}</div>
            </div>
          ))}
        </div>

        {/* Center column — Video circle */}
        <div
          className="flex items-center justify-center order-first lg:order-none"
          style={{ alignSelf: 'center' }}
        >
          <div
            style={{
              position: 'relative',
              borderRadius: '50%',
              overflow: 'hidden',
              width: 'clamp(200px, 22vw, 400px)',
              height: 'clamp(200px, 22vw, 400px)',
              flexShrink: 0,
            }}
          >
            <HlsVideo />
          </div>
        </div>

        {/* Right column — Positives */}
        <div
          className="flex flex-col"
          style={{
            gap: '12px',
            fontSize: 'clamp(13px, 1.15vw, 17px)',
          }}
        >
          {positives.map((text, i) => (
            <div
              key={i}
              className="flex flex-col"
              style={{
                gap: '12px',
                padding: 'clamp(12px, 0.97vw, 16px) clamp(14px, 1.25vw, 20px)',
                borderRadius: '18px',
                backgroundColor: 'rgb(255, 255, 255)',
                boxShadow: '0 3px 9.1px #3f4a7e0d, 0 1px 29px #3f4a7e1a',
              }}
            >
              <img
                src={CHECK_ICON}
                alt=""
                aria-hidden
                style={{ width: 'clamp(16px, 1.25vw, 20px)', flexShrink: 0 }}
              />
              <div style={{ color: 'rgb(26, 11, 84)' }}>{text}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
};

export default FreedomSection;
