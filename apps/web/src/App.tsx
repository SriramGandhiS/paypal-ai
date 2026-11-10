import React, { useState, useEffect } from 'react';
import { Navbar } from './components/Navbar';
import { DashboardScreen } from './screens/DashboardScreen';
import { CaseDetailScreen } from './screens/CaseDetailScreen';
import { PassportScreen } from './screens/PassportScreen';
import { TruthGapGraphScreen } from './screens/TruthGapGraphScreen';
import { RecoveryMirrorScreen } from './screens/RecoveryMirrorScreen';
import { GoldenPathDemoScreen } from './screens/GoldenPathDemoScreen';
import { api } from './api/client';
import { ShieldCheck, Info } from 'lucide-react';
import {
  LandingScreen,
  AuthScreen,
  NewPaymentScreen,
  TimelineScreen,
  SenderLensScreen,
  ReceiverLensScreen,
  SafeActionScreen,
  EvidenceScreen,
  NotificationsScreen,
  OperationsScreen,
  SettingsScreen,
} from './features/screens/Screens';

export const App: React.FC = () => {
  const [currentTab, setCurrentTab] = useState<string>('dashboard');
  const [selectedCaseId, setSelectedCaseId] = useState<string | null>(null);
  const [isSimulated, setIsSimulated] = useState<boolean>(api.isSimulationMode());
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  useEffect(() => {
    // Initial backend ping
    api.checkBackendLiveness().then((live) => {
      setIsSimulated(!live);
    });
  }, []);

  const showToast = (msg: string) => {
    setToastMessage(msg);
    setTimeout(() => setToastMessage(null), 3500);
  };

  const handleToggleMode = () => {
    const nextMode = !isSimulated;
    api.setSimulationMode(nextMode);
    setIsSimulated(nextMode);
    showToast(nextMode ? 'Switched to Deterministic Sandbox Simulator' : 'Switched to Live Backend API');
  };

  const handleSelectCase = (id: string) => {
    setSelectedCaseId(id);
    setCurrentTab('cases');
  };

  const handleTabChange = (tab: string) => {
    setCurrentTab(tab);
    if (tab !== 'cases') {
      setSelectedCaseId(null);
    } else if (!selectedCaseId) {
      // Default to primary demo case if opening cases tab directly
      setSelectedCaseId('c101-demo-4567-89ab-cdef00000001');
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-indigo-500 selection:text-white">
      {/* Top Navigation */}
      <Navbar
        currentTab={currentTab}
        onTabChange={handleTabChange}
        isSimulated={isSimulated}
        onToggleMode={handleToggleMode}
      />

      {/* Floating Notification Toast */}
      {toastMessage && (
        <div className="fixed bottom-6 right-6 z-50 animate-slide-up flex items-center gap-2.5 px-4 py-2.5 rounded-xl bg-slate-900 border border-indigo-500/40 text-slate-200 text-xs shadow-2xl shadow-black/80 font-medium">
          <Info className="w-4 h-4 text-cyan-400 shrink-0" />
          <span>{toastMessage}</span>
        </div>
      )}

      {/* Main Content Body */}
      <main className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 pt-8">
        {currentTab === 'dashboard' && (
          <DashboardScreen
            onSelectCase={handleSelectCase}
            onNavigate={handleTabChange}
          />
        )}

        {currentTab === 'cases' && (
          selectedCaseId ? (
            <CaseDetailScreen
              caseId={selectedCaseId}
              onBack={() => {
                setSelectedCaseId(null);
                setCurrentTab('dashboard');
              }}
              onNavigate={handleTabChange}
            />
          ) : (
            <DashboardScreen
              onSelectCase={handleSelectCase}
              onNavigate={handleTabChange}
            />
          )
        )}

        {currentTab === 'demo' && (
          <GoldenPathDemoScreen
            onNavigate={handleTabChange}
            onSelectCase={handleSelectCase}
          />
        )}

        {currentTab === 'passport' && (
          <PassportScreen />
        )}

        {currentTab === 'graph' && (
          <TruthGapGraphScreen />
        )}

        {currentTab === 'mirror' && (
          <RecoveryMirrorScreen />
        )}

        {/* Vijay Extended 10 Screens */}
        {currentTab === 'landing' && (
          <LandingScreen onNavigate={() => handleTabChange('new-payment')} />
        )}

        {currentTab === 'auth' && (
          <AuthScreen onLogin={() => handleTabChange('dashboard')} />
        )}

        {currentTab === 'new-payment' && (
          <NewPaymentScreen onSimulateTimeout={() => handleTabChange('cases')} />
        )}

        {currentTab === 'timeline' && (
          <TimelineScreen />
        )}

        {currentTab === 'sender-lens' && (
          <SenderLensScreen />
        )}

        {currentTab === 'receiver-lens' && (
          <ReceiverLensScreen />
        )}

        {currentTab === 'safe-action' && (
          <SafeActionScreen />
        )}

        {currentTab === 'evidence' && (
          <EvidenceScreen />
        )}

        {currentTab === 'notifications' && (
          <NotificationsScreen />
        )}

        {currentTab === 'operations' && (
          <OperationsScreen />
        )}

        {currentTab === 'settings' && (
          <SettingsScreen />
        )}
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-900 bg-slate-950/90 py-6 mt-16 text-xs text-slate-500">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <ShieldCheck className="w-4 h-4 text-indigo-400" />
            <span className="font-semibold text-slate-400">PayBridge AI</span>
            <span>— Financial Truth & Reconciliation Architecture</span>
          </div>
          <div className="flex items-center gap-4 text-slate-400 font-mono text-[11px]">
            <span>Provider: PayPal Sandbox</span>
            <span>•</span>
            <span>Audit: Tier 1–6 SHA-256</span>
            <span>•</span>
            <span>AI Mutations: Prohibited</span>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default App;
