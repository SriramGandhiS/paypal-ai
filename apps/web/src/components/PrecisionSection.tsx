import React from 'react';

const LOGO_ICON =
  'https://cdn.prod.website-files.com/6720dd1ab6df0da205830ab1/6870f623cf3df417ce45df05_icon%20logo%20eternacloud.png';

const LINE_GRADIENT =
  'linear-gradient(rgb(28, 78, 255), rgb(254, 136, 27) 0%, rgb(172, 36, 255) 25%, rgb(247, 159, 255) 50%, rgb(255, 214, 0) 66%, rgb(254, 136, 27) 84%, rgba(254, 136, 27, 0) 102%)';

export interface PillarData {
  label: string;
  items: string[];
  leftVw: number;
  bottomVw: number;
}

// Custom words tailored for PayBridge AI cross-border truth architecture
const PILLARS: PillarData[] = [
  { label: 'Scopes',     items: ['conditions', 'capacity', 'specs', 'timelines'],     leftVw: 2.8,  bottomVw: 7     },
  { label: 'Integrates', items: ['civil', 'mechanical', 'electrical', 'controls'],     leftVw: 22.4, bottomVw: 9.08  },
  { label: 'Certifies',  items: ['redundancy', 'testing', 'compliance', 'sign-offs'], leftVw: 41.2, bottomVw: 11.16 },
  { label: 'Activates',  items: ['cutover', 'runbooks', 'handoff', 'SLAs'],           leftVw: 61.1, bottomVw: 13.24 },
];

export interface PrecisionSectionProps {
  pillars?: PillarData[];
  badgeText?: string;
  headingLine1?: string;
  headingLine2?: string;
  subtext?: string;
}

export const PrecisionSection: React.FC<PrecisionSectionProps> = ({
  pillars = PILLARS,
  badgeText = 'Structured Delivery',
  headingLine1 = 'One integrated, end-to-end system.',
  headingLine2 = 'Compounding operational value.',
  subtext = 'PayBridge teams capture, align, validate and deliver exactly what keeps your programs on track.',
}) => {
  return (
    <section
      style={{
        backgroundImage:
          'url("https://images.higgs.ai/?default=1&output=webp&url=https%3A%2F%2Fd8j0ntlcm91z4.cloudfront.net%2Fuser_38xzZboKViGWJOttwIXH07lWA1P%2Fhf_20260418_125638_553b96dc-a1fd-4b2b-81a9-ed7daa80006e.png&w=1280&q=85")',
        backgroundSize: 'cover',
        backgroundPosition: 'center',
        backgroundRepeat: 'no-repeat',
        width: '100%',
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        textAlign: 'center',
        padding: 'clamp(48px, 8vw, 120px) clamp(16px, 4vw, 60px) clamp(48px, 5.56vw, 80px)',
        gap: 'clamp(32px, 4vw, 56px)',
      }}
    >
      {/* Block 1 — Header */}
      <div
        style={{
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
          gap: '36px',
        }}
      >
        {/* Badge pill */}
        <div
          style={{
            backgroundColor: 'rgb(249, 249, 249)',
            display: 'flex',
            alignItems: 'center',
            gap: '8px',
            fontSize: 'clamp(14px, 1.1vw, 18px)',
            fontWeight: 500,
            borderRadius: '36px',
            padding: 'clamp(8px, 0.9vw, 14px) clamp(12px, 1.25vw, 20px)',
            color: 'rgb(26, 11, 84)',
            whiteSpace: 'nowrap',
          }}
        >
          <svg
            width="19"
            height="18"
            style={{ flexShrink: 0 }}
            viewBox="0 0 17 16"
            fill="none"
          >
            <g clipPath="url(#prec-clip)">
              <circle cx="8.5" cy="8" r="7" stroke="#c86fff" fill="none" />
              <path
                d="M9.5 11.5V10.5H7.5V11.5H9.5ZM7.5 14.5C7.5 15.0523 7.94772 15.5 8.5 15.5C9.05228 15.5 9.5 15.0523 9.5 14.5H7.5ZM8.5 11.5H7.5V14.5H8.5H9.5V11.5H8.5Z"
                fill="rgb(200, 111, 255)"
              />
              <path
                d="M12 7H11V9H12V7ZM15 9C15.5523 9 16 8.55228 16 8C16 7.44772 15.5523 7 15 7V9ZM12 8V9H15V8V7L12 7V8Z"
                fill="rgb(200, 111, 255)"
              />
              <path
                d="M5 9H6V7H5V9ZM2 7C1.44772 7 1 7.44772 1 8C1 8.55228 1.44772 9 2 9V7ZM5 8V7H2V8V9H5V8Z"
                fill="rgb(200, 111, 255)"
              />
              <path
                d="M7.5 4.5V5.5H9.5V4.5H7.5ZM9.5 1.5C9.5 0.947715 9.05228 0.5 8.5 0.5C7.94772 0.5 7.5 0.947715 7.5 1.5H9.5ZM8.5 4.5H9.5V1.5H8.5H7.5V4.5H8.5Z"
                fill="rgb(200, 111, 255)"
              />
            </g>
            <defs>
              <clipPath id="prec-clip">
                <rect width="16" height="16" fill="white" transform="translate(0.5)" />
              </clipPath>
            </defs>
          </svg>
          {badgeText}
        </div>

        {/* Heading + subtext */}
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            maxWidth: 'clamp(700px, 60vw, 900px)',
            gap: '22px',
          }}
        >
          <h2
            style={{
              fontSize: 'clamp(28px, 4vw, 56px)',
              fontWeight: 500,
              color: 'rgb(26, 11, 84)',
              lineHeight: 1.15,
              margin: 0,
            }}
          >
            <span className="sm:whitespace-nowrap" style={{ display: 'block' }}>
              {headingLine1}
            </span>
            <span
              style={{
                backgroundImage:
                  'linear-gradient(90deg, rgb(43, 167, 255), rgb(202, 69, 255) 50%, rgb(254, 136, 27))',
                WebkitBackgroundClip: 'text',
                backgroundClip: 'text',
                WebkitTextFillColor: 'transparent',
                color: 'transparent',
                paddingBottom: '0.3vw',
                display: 'block',
              }}
            >
              {headingLine2}
            </span>
          </h2>

          <p
            style={{
              fontSize: 'clamp(15px, 1.2vw, 20px)',
              color: 'rgb(169, 151, 206)',
              margin: 0,
            }}
          >
            {subtext}
          </p>
        </div>
      </div>

      {/* Block 2 — Pillars container */}
      <div
        style={{
          width: '100%',
          maxWidth: '82.292vw',
          margin: '0 auto',
        }}
      >
        {/* Desktop pillars — hidden sm:block */}
        <div
          className="hidden sm:block"
          style={{
            position: 'relative',
            width: '82.292vw',
            height: '31.94vw',
            color: 'rgb(26, 11, 84)',
          }}
        >
          {pillars.map((pillar, index) => (
            <div
              key={index}
              style={{
                position: 'absolute',
                bottom: `${pillar.bottomVw}vw`,
                left: `${pillar.leftVw}vw`,
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'center',
                justifyContent: 'flex-start',
              }}
            >
              {/* Chip */}
              <div
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  backgroundImage:
                    'linear-gradient(135deg, rgb(255, 255, 255), rgba(255, 255, 255, 0.6))',
                  fontSize: '18px',
                  fontWeight: 500,
                  borderRadius: '20px',
                  paddingTop: '0.972vw',
                  paddingBottom: '0.972vw',
                  paddingLeft: '1.736vw',
                  paddingRight: '1.736vw',
                  whiteSpace: 'nowrap',
                  gap: '8px',
                }}
              >
                <img
                  src={LOGO_ICON}
                  alt=""
                  style={{ width: '1.111vw', height: 'auto', display: 'inline-block' }}
                />
                {pillar.label}
              </div>

              {/* Line + items wrapper */}
              <div
                style={{
                  position: 'relative',
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: 'center',
                  justifyContent: 'flex-end',
                }}
              >
                {/* Items container */}
                <div
                  style={{
                    position: 'absolute',
                    top: '0.56vw',
                    left: '1.94vw',
                    display: 'flex',
                    flexDirection: 'column',
                    gap: '4px',
                    fontSize: '16px',
                    alignItems: 'flex-start',
                    justifyContent: 'space-between',
                  }}
                >
                  {pillar.items.map((item, itemIdx) => (
                    <div
                      key={itemIdx}
                      style={{
                        paddingTop: '0.69vw',
                        paddingBottom: '0.69vw',
                        paddingLeft: '1.04vw',
                        paddingRight: '1.04vw',
                        display: 'flex',
                        alignItems: 'flex-start',
                        whiteSpace: 'nowrap',
                      }}
                    >
                      {item}
                    </div>
                  ))}
                </div>

                {/* Vertical gradient line */}
                <div
                  style={{
                    backgroundImage: LINE_GRADIENT,
                    width: '1px',
                    height: '14.24vw',
                  }}
                />
              </div>
            </div>
          ))}
        </div>

        {/* Mobile pillars — flex flex-col sm:hidden w-full */}
        <div
          className="flex flex-col sm:hidden w-full"
          style={{
            color: 'rgb(26, 11, 84)',
            gap: 0,
          }}
        >
          {pillars.map((pillar, index) => {
            const isRight = index % 2 !== 0;
            return (
              <div
                key={index}
                style={{
                  display: 'flex',
                  flexDirection: 'column',
                  alignItems: isRight ? 'flex-end' : 'flex-start',
                  width: '100%',
                  paddingBottom: '8px',
                }}
              >
                {/* Chip */}
                <div
                  style={{
                    display: 'inline-flex',
                    alignItems: 'center',
                    backgroundImage:
                      'linear-gradient(135deg, rgb(255, 255, 255), rgba(255, 255, 255, 0.6))',
                    fontSize: '15px',
                    fontWeight: 500,
                    borderRadius: '20px',
                    padding: '10px 18px',
                    whiteSpace: 'nowrap',
                    gap: '7px',
                  }}
                >
                  <img src={LOGO_ICON} alt="" style={{ width: 16, height: 'auto' }} />
                  {pillar.label}
                </div>

                {/* Line + items row */}
                <div
                  style={{
                    display: 'flex',
                    flexDirection: isRight ? 'row-reverse' : 'row',
                    alignItems: 'stretch',
                    width: '100%',
                  }}
                >
                  {/* Vertical line */}
                  <div
                    style={{
                      width: '1px',
                      flexShrink: 0,
                      backgroundImage: LINE_GRADIENT,
                      marginLeft: isRight ? 0 : '22px',
                      marginRight: isRight ? '22px' : 0,
                      minHeight: '120px',
                    }}
                  />

                  {/* Items */}
                  <div
                    style={{
                      display: 'flex',
                      flexDirection: 'column',
                      gap: 0,
                      paddingLeft: isRight ? 0 : '20px',
                      paddingRight: isRight ? '20px' : 0,
                      paddingTop: '8px',
                      paddingBottom: '8px',
                      alignItems: isRight ? 'flex-end' : 'flex-start',
                    }}
                  >
                    {pillar.items.map((item, itemIdx) => (
                      <div
                        key={itemIdx}
                        style={{
                          fontSize: '14px',
                          color: 'rgb(100, 80, 160)',
                          padding: '8px 0',
                        }}
                      >
                        {item}
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
};

export default PrecisionSection;
