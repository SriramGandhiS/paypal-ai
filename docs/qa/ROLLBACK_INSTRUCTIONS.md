# PayBridge AI — Golden Path Rollback & Safe Recovery Instructions

**Specification Reference:** Section 13  
**Active Branch:** `feature/golden-path-paypal-flow`  
**Base Commit:** `8c9b6d5`  

---

## 1. Zero-Risk Rollback Procedure

If the user does not approve the Golden Path redesign, execute the following commands to restore the baseline:

### Option A: Complete Branch Rollback (Instant)
```powershell
# 1. Switch back to previous feature branch or master
git checkout feature/enterprise-investigation-suite

# 2. Rebuild the frontend
npm --prefix "D:\Projects\PayBridgeAI\apps\web" run build

# 3. Verify server status on port 3000
python -c "import urllib.request; print(urllib.request.urlopen('http://localhost:3000/').status)"
```

### Option B: Component-Only Reversion
To revert only the Golden Path files while preserving other work on this branch:
```powershell
# 1. Restore baseline Golden Path React component
git checkout 8c9b6d5 -- apps/web/src/screens/GoldenPathDemoScreen.tsx

# 2. Restore baseline HTML
git checkout 8c9b6d5 -- apps/web/index.html

# 3. Recompile production bundle
npm --prefix "D:\Projects\PayBridgeAI\apps\web" run build
```

---

## 2. Protected Scopes
The following files and configurations were **never modified** and remain 100% intact:
- Website background, gradients, wave graphics, and decorative canvas.
- Global navigation bar, header, and conveyor belt cards.
- Landing page, Dashboard, Cases, Reports, and Passport screens.
- Global CSS, design tokens, and Tailwind configuration.
- Go backend microservices, DB schemas, and Dinny D1–D7 enterprise modules.
