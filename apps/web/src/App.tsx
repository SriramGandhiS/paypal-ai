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
  const [currentTab, setCurrentTab] = useState<string>('landing');
  const [selectedCaseId, setSelectedCaseId] = useState<string | null>(null);
  const [isSimulated, setIsSimulated] = useState<boolean>(api.isSimulationMode());
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  useEffect(() => {
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
      setSelectedCaseId('c101-demo-4567-89ab-cdef00000001');
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-indigo-500 selection:text-white">
      {/* Top Navigation Bar with 16 Screens Selector */}
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

      {/* Main Content: Full-Bleed for Landing Showcase (Distill Scanner), Padded Container for Enterprise Views */}
      <main className={`flex-1 w-full mx-auto ${currentTab === 'landing' ? 'max-w-none p-0' : 'max-w-7xl px-4 sm:px-6 lg:px-8 pt-8'}`}>
        {/* Screen 1: Landing Showcase (Distill Scanner Hero + Invariants) */}
        {currentTab === 'landing' && (
          <LandingScreen onNavigate={handleTabChange} />
        )}

        {/* Screen 3: Dashboard Matrix */}
        {currentTab === 'dashboard' && (
          <DashboardScreen
            onSelectCase={handleSelectCase}
            onNavigate={handleTabChange}
          />
        )}

        {/* Screen 5: Payment Cases & Detail View */}
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

        {/* Screen 16: Golden Path Stepper */}
        {currentTab === 'demo' && (
          <GoldenPathDemoScreen
            onNavigate={handleTabChange}
            onSelectCase={handleSelectCase}
          />
        )}

        {/* Screen 11: Truth Passport */}
        {currentTab === 'passport' && (
          <PassportScreen />
        )}

        {/* Screen 7: Truth Gap Graph */}
        {currentTab === 'graph' && (
          <TruthGapGraphScreen />
        )}

        {/* Screen 11 (Signature): Recovery Mirror */}
        {currentTab === 'mirror' && (
          <RecoveryMirrorScreen />
        )}

        {/* Screen 2: Login & RBAC Auth */}
        {currentTab === 'auth' && (
          <AuthScreen onLogin={(role) => {
            handleTabChange('dashboard');
            showToast(`Authenticated as ${role.toUpperCase()}`);
          }} />
        )}

        {/* Screen 4: Checkout & Timeout Hero */}
        {currentTab === 'new-payment' && (
          <NewPaymentScreen onSimulateTimeout={() => {
            handleTabChange('cases');
            showToast('Gateway Timeout 504 Simulated: DO NOT REPEAT CHECKOUT lock engaged');
          }} />
        )}

        {/* Screen 6: Lifecycle Timeline */}
        {currentTab === 'timeline' && (
          <TimelineScreen />
        )}

        {/* Screen 8: Sender Two-Lens View */}
        {currentTab === 'sender-lens' && (
          <SenderLensScreen />
        )}

        {/* Screen 9: Receiver Two-Lens View */}
        {currentTab === 'receiver-lens' && (
          <ReceiverLensScreen />
        )}

        {/* Screen 10: Next Safe Action Engine */}
        {currentTab === 'safe-action' && (
          <SafeActionScreen />
        )}

        {/* Screen 12: Evidence Explorer */}
        {currentTab === 'evidence' && (
          <EvidenceScreen />
        )}

        {/* Screen 13: Notification Center */}
        {currentTab === 'notifications' && (
          <NotificationsScreen />
        )}

        {/* Screen 14: Operations & Webhook Ingestion */}
        {currentTab === 'operations' && (
          <OperationsScreen />
        )}

        {/* Screen 15: Privacy & Consent Settings */}
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
