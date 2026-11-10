import React from 'react';
import { 
  ShieldCheck, 
  LayoutDashboard, 
  FileText, 
  PlayCircle, 
  Award, 
  Network, 
  GitFork, 
  Radio
} from 'lucide-react';

interface NavbarProps {
  currentTab: string;
  onTabChange: (tab: string) => void;
  isSimulated: boolean;
  onToggleMode: () => void;
}

export const Navbar: React.FC<NavbarProps> = ({
  currentTab,
  onTabChange,
  isSimulated,
  onToggleMode,
}) => {
  const tabs = [
    { id: 'dashboard', label: 'Dashboard', icon: LayoutDashboard },
    { id: 'cases', label: 'Payment Cases', icon: FileText },
    { id: 'demo', label: 'Golden Path Demo', icon: PlayCircle },
    { id: 'passport', label: 'Truth Passport', icon: Award },
    { id: 'graph', label: 'Truth Gap Graph', icon: Network },
    { id: 'mirror', label: 'Recovery Mirror', icon: GitFork },
  ];

  return (
    <header className="sticky top-0 z-50 bg-slate-950/80 backdrop-blur-md border-b border-slate-800">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Brand Logo */}
          <div className="flex items-center gap-3 cursor-pointer" onClick={() => onTabChange('dashboard')}>
            <div className="w-10 h-10 rounded-xl bg-gradient-to-tr from-indigo-600 via-indigo-500 to-cyan-400 p-0.5 flex items-center justify-center shadow-lg shadow-indigo-500/20">
              <div className="w-full h-full bg-slate-950 rounded-[10px] flex items-center justify-center">
                <ShieldCheck className="w-5 h-5 text-cyan-400" />
              </div>
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="font-bold text-lg text-slate-100 tracking-tight">PayBridge <span className="text-cyan-400 font-extrabold">AI</span></span>
                <span className="text-[10px] px-2 py-0.5 rounded bg-indigo-500/20 text-indigo-300 border border-indigo-500/30 font-mono font-medium">MVP</span>
              </div>
              <p className="text-xs text-slate-400 font-medium">Financial Truth & Reconciliation</p>
            </div>
          </div>

          {/* Navigation Links */}
          <nav className="hidden md:flex items-center gap-1">
            {tabs.map((tab) => {
              const Icon = tab.icon;
              const active = currentTab === tab.id;
              return (
                <button
                  key={tab.id}
                  onClick={() => onTabChange(tab.id)}
                  className={`flex items-center gap-2 px-3.5 py-2 rounded-lg text-sm font-medium transition-all ${
                    active 
                      ? 'bg-indigo-600/20 text-indigo-300 border border-indigo-500/30 shadow-sm'
                      : 'text-slate-400 hover:text-slate-200 hover:bg-slate-900/60'
                  }`}
                >
                  <Icon className={`w-4 h-4 ${active ? 'text-indigo-400' : 'text-slate-400'}`} />
                  {tab.label}
                </button>
              );
            })}
          </nav>

          {/* Right Mode Switcher */}
          <div className="flex items-center gap-3">
            <button
              onClick={onToggleMode}
              title="Click to toggle between Live Backend and Sandbox Simulator mode"
              className={`flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-mono font-medium border transition-all ${
                isSimulated
                  ? 'bg-cyan-950/40 text-cyan-300 border-cyan-800/60 hover:bg-cyan-900/40'
                  : 'bg-emerald-950/40 text-emerald-300 border-emerald-800/60 hover:bg-emerald-900/40'
              }`}
            >
              <Radio className={`w-3.5 h-3.5 animate-pulse ${isSimulated ? 'text-cyan-400' : 'text-emerald-400'}`} />
              <span>{isSimulated ? 'SIMULATOR MODE' : 'LIVE API CONNECTED'}</span>
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};
