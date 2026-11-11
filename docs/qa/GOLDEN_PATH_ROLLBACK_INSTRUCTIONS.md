# Golden Path Demo — Rollback & Safe Reversion Instructions

**Repository:** `D:\Projects\PayBridgeAI`  
**Dedicated Branch:** `feature/golden-path-redesign`  
**Baseline Git Checkpoint:** `921d7916c6c80d8bb257b4dc4a6361d87f888433`  

---

## 1. Golden Path Safe Rollback Policy

In strict accordance with **Section 9: Rollback and User Approval — Mandatory**:
- The user has **final design authority**.
- If the user indicates that the redesigned Golden Path is not satisfactory or chooses to revert, the original design can be restored completely within seconds without affecting any unrelated work or other website components.

---

## 2. Immediate One-Step Rollback Command

To instantly revert all Golden Path changes back to the exact pre-redesign baseline:

```powershell
# In D:\Projects\PayBridgeAI
git checkout feature/frontend-complete
```

Or to restore the files on the current branch without switching branches:

```powershell
git checkout 921d7916c6c80d8bb257b4dc4a6361d87f888433 -- apps/web/index.html apps/web/src/screens/GoldenPathDemoScreen.tsx
git commit -m "revert: restore original golden path design from baseline checkpoint"
```

---

## 3. Pre-Redesign Baseline Evidence

- **Pre-Redesign Screenshot:** `docs/qa/golden_path_before_redesign.png`
- **Pre-Redesign Audit:** `docs/qa/GOLDEN_PATH_BEFORE_REDESIGN.md`
- **Baseline Commit:** `921d791` (`fix(conveyor): eliminate card text overlays and adjust typography scale`)
