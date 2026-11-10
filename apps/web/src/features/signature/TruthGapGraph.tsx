import React, { useState } from 'react';
import { TruthGapNode, TruthGapEdge, TruthGapItem } from '../../types';

interface TruthGapGraphProps {
  nodes?: TruthGapNode[];
  edges?: TruthGapEdge[];
  gaps?: TruthGapItem[];
}

export const TruthGapGraph: React.FC<TruthGapGraphProps> = ({
  nodes = defaultNodes,
  edges = defaultEdges,
  gaps = defaultGaps,
}) => {
  const [selectedNode, setSelectedNode] = useState<TruthGapNode | null>(null);

  const lanes = ['Browser', 'Merchant DB', 'Provider', 'Recipient'];
  const laneHeight = 80;
  const graphWidth = 720;
  const graphHeight = lanes.length * laneHeight + 40;

  const getLaneY = (lane: string) => {
    const idx = lanes.indexOf(lane);
    return idx >= 0 ? idx * laneHeight + 50 : 50;
  };

  const getNodeX = (nodeId: string) => {
    switch (nodeId) {
      case 'b1': return 120;
      case 'm1': return 260;
      case 'p1': return 420;
      case 'p2': return 580;
      case 'r1': return 580;
      default: return 300;
    }
  };

  return (
    <div className="card" style={{ padding: '24px', background: 'var(--surface-1)' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '16px' }}>
        <div>
          <h3 style={{ fontSize: '18px', fontWeight: 600, color: 'var(--text-1)' }}>
            Truth Gap Graph — Multi-System Reality Matrix
          </h3>
          <p style={{ fontSize: '13px', color: 'var(--text-2)' }}>
            Observes independent systems without collapsing contradictory states into a false single truth.
          </p>
        </div>
        <div style={{ display: 'flex', gap: '8px' }}>
          <span className="badge badge-unknown">Contradiction Detected</span>
          <span className="badge badge-verified">Provider Primary</span>
        </div>
      </div>

      <div style={{ overflowX: 'auto' }}>
        <svg
          viewBox={`0 0 ${graphWidth} ${graphHeight}`}
          style={{ width: '100%', minWidth: '680px', height: 'auto', background: 'var(--surface-2)', borderRadius: '8px', border: '1px solid var(--border)' }}
        >
          {/* Lane dividers and labels */}
          {lanes.map((lane, idx) => (
            <g key={lane}>
              <line
                x1="0"
                y1={idx * laneHeight + 30}
                x2={graphWidth}
                y2={idx * laneHeight + 30}
                stroke="var(--border)"
                strokeDasharray="4 4"
              />
              <rect
                x="10"
                y={idx * laneHeight + 35}
                width="100"
                height="22"
                fill="var(--surface-sunken)"
                rx="4"
              />
              <text
                x="16"
                y={idx * laneHeight + 50}
                fontSize="11"
                fontWeight="600"
                fill="var(--text-2)"
              >
                {lane.toUpperCase()}
              </text>
            </g>
          ))}

          {/* Relationship Edges */}
          {edges.map((edge, i) => {
            const fromNode = nodes.find(n => n.id === edge.from);
            const toNode = nodes.find(n => n.id === edge.to);
            if (!fromNode || !toNode) return null;

            const x1 = getNodeX(fromNode.id) + 60;
            const y1 = getLaneY(fromNode.lane);
            const x2 = getNodeX(toNode.id) - 60;
            const y2 = getLaneY(toNode.lane);

            const isContradiction = edge.relation === 'CONTRADICTS';

            return (
              <g key={i}>
                <path
                  d={`M ${x1} ${y1} C ${(x1 + x2) / 2} ${y1}, ${(x1 + x2) / 2} ${y2}, ${x2} ${y2}`}
                  fill="none"
                  stroke={isContradiction ? 'var(--unknown)' : 'var(--accent)'}
                  strokeWidth={isContradiction ? '2.5' : '1.5'}
                  strokeDasharray={isContradiction ? '6 4' : 'none'}
                />
                {isContradiction && (
                  <circle cx={(x1 + x2) / 2} cy={(y1 + y2) / 2} r="9" fill="var(--unknown)" />
                )}
              </g>
            );
          })}

          {/* Nodes */}
          {nodes.map(node => {
            const cx = getNodeX(node.id);
            const cy = getLaneY(node.lane);
            const isSelected = selectedNode?.id === node.id;

            return (
              <g
                key={node.id}
                onClick={() => setSelectedNode(node)}
                style={{ cursor: 'pointer' }}
              >
                <rect
                  x={cx - 65}
                  y={cy - 22}
                  width="130"
                  height="44"
                  rx="6"
                  fill={node.isAmbiguous ? 'var(--unknown-bg)' : 'var(--surface-1)'}
                  stroke={node.isAmbiguous ? 'var(--unknown)' : isSelected ? 'var(--accent)' : 'var(--border-strong)'}
                  strokeWidth={isSelected || node.isAmbiguous ? 2 : 1}
                />
                <text
                  x={cx}
                  y={cy - 4}
                  textAnchor="middle"
                  fontSize="12"
                  fontWeight="600"
                  fill="var(--text-1)"
                >
                  {node.label}
                </text>
                <text
                  x={cx}
                  y={cy + 12}
                  textAnchor="middle"
                  fontSize="10"
                  fontWeight="bold"
                  fill={node.isAmbiguous ? 'var(--unknown)' : 'var(--verified)'}
                >
                  {node.status} (E{node.tier})
                </text>
              </g>
            );
          })}
        </svg>
      </div>

      {/* Contradiction details drawer / summary */}
      <div style={{ marginTop: '16px', display: 'flex', flexDirection: 'column', gap: '8px' }}>
        {gaps.map(gap => (
          <div
            key={gap.id}
            style={{
              padding: '12px 16px',
              backgroundColor: 'var(--unknown-bg)',
              border: '1px solid var(--unknown)',
              borderRadius: '6px',
              fontSize: '13px',
              color: 'var(--text-1)',
            }}
          >
            <div style={{ display: 'flex', alignItems: 'center', gap: '8px', marginBottom: '4px' }}>
              <span className="badge badge-unknown">{gap.severity}</span>
              <strong>{gap.title}</strong>
            </div>
            <p style={{ color: 'var(--text-2)', margin: 0 }}>{gap.description}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

const defaultNodes: TruthGapNode[] = [
  { id: 'b1', lane: 'Browser', label: 'HTTP Checkout', status: 'API_TIMEOUT', tier: 3, isAmbiguous: true },
  { id: 'm1', lane: 'Merchant DB', label: 'Order #PB-8821', status: 'WAITING_CAPTURE', tier: 4 },
  { id: 'p1', lane: 'Provider', label: 'PayPal Order', status: 'APPROVED', tier: 1 },
  { id: 'p2', lane: 'Provider', label: 'PayPal Capture', status: 'COMPLETED', tier: 1 },
  { id: 'r1', lane: 'Recipient', label: 'Bank Settlement', status: 'NOT_OBSERVABLE', tier: 5 },
];

const defaultEdges: TruthGapEdge[] = [
  { from: 'b1', to: 'm1', relation: 'PRECEDES' },
  { from: 'b1', to: 'p2', relation: 'CONTRADICTS' },
  { from: 'p1', to: 'p2', relation: 'PRECEDES' },
];

const defaultGaps: TruthGapItem[] = [
  {
    id: 'GAP-01',
    title: 'Gap G02: Browser Timeout vs Authoritative Provider Capture',
    severity: 'ATTENTION',
    description: 'The browser client timed out after 15s and marked the request unknown, whereas PayPal processed and completed the capture. PayBridge reconciled the original capture without issuing a duplicate transaction.',
  },
];
