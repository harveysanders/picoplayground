# Sulfur Gas Sensing Architecture (High-Level)

## Scope

This document summarizes:

- Recommended **sensor variants** for SO₂ and H₂S
- Required **AFE (Analog Front End) architecture**
- **Interconnect** (connectors/cabling)
  for urban / near-industrial outdoor air-quality monitoring.

---

## Recommended Sensor Choices

### Sensor family

- **Alphasense A-series, 4-electrode electrochemical sensors**
- [SO₂ sensors](https://store.ametekmocon.com/alphasense/gas-sensors/target-gases/sulphur-dioxide/)
- [H₂S sensors](https://store.ametekmocon.com/alphasense/gas-sensors/target-gases/hydrogen-sulphide/)

Reason:

- Designed for low-level (ppb → low-ppm) detection
- Better baseline stability and drift behavior than B-series
- Appropriate for outdoor environmental monitoring

---

### Variant selection

#### **A4**

- Standard 4-electrode can
- Best general choice for prototyping and early deployments
- Widely used and well characterized

#### **A4+**

- Same as A4, with improved long-term stability and reduced drift
- Same pinout and electrical behavior
- **Preferred for unattended or long-term outdoor deployments**

#### **A41**

- Same electrochemistry as A4/A4+
- Lower-profile / packaging-optimized mechanical form
- Electrically interchangeable with A4/A4+
- Useful if enclosure height or airflow geometry is constrained

**Recommendation**

- Default: **A4+**
- Use **A41** only if mechanical layout requires it

---

### Gas-specific recommendations

- **SO₂:** SO2-A4+ (or SO2-A41+ if packaging requires)
- **H₂S:** H2S-A4+ (or H2S-A41+ if packaging requires)

---

## AFE (Analog Front End) Architecture

### Key constraint

Alphasense **A-series AFEs have fixed analog “slots”** based on gas families.
Reducing gases (CO / SO₂ / H₂S) **share the same bias/gain family**.

As a result:

- A single 2-way AFE **cannot** support SO₂ and H₂S simultaneously
- Each sulfur sensor requires its **own reducing-gas slot**

---

### Recommended AFE approach

Use **two identical 2-way AFEs**, each with:

- **NO₂ (oxidizing gas) slot**
- **Generic reducing-gas slot (CO / SO₂ / H₂S)**

Populate as follows:

- **AFE #1**

  - Reducing slot → **SO₂ sensor**
  - Oxidizing slot → **NO₂ sensor**

- **AFE #2**
  - Reducing slot → **H₂S sensor**
  - Oxidizing slot → **NO₂ sensor**

Rationale:

- Allows SO₂ and H₂S to operate simultaneously
- NO₂ provides valuable atmospheric context and interference insight
- Avoids unsupported dual-sulfur configurations on 2-way AFEs
- Keeps architecture simple and within Alphasense’s intended use

---

### Notes

- Leaving the second channel populated with NO₂ is intentional
- Leaving channels unpopulated is acceptable, but provides less context
- 3-way AFEs with two reducing-gas slots are an alternative, but larger and higher cost

---

## Interconnect / Cabling

### Connector type

- **10-pin IDC ribbon connector**
- **2.54 mm (0.1") pitch**
- **2×5 female IDC socket**

This is the standard interface used on most Alphasense A-series AFEs.

---

### Recommended setup

- 10-pin IDC ribbon cables
- 10-pin IDC breakout boards or mating 2×5 headers on the main PCB
- Clear pin-1 orientation marking

Always verify the **exact AFE SKU pinout** before final PCB layout.

---

## Summary

- Sensors: **SO₂-A4+ and H₂S-A4+** (A41 variants only if mechanically required)
- AFEs: **2 × 2-way AFEs**, each with **NO₂ + one sulfur sensor**
- Interconnect: **10-pin, 2.54 mm IDC ribbon cables and breakouts**

This architecture fits current A-series AFE constraints, supports simultaneous SO₂ and H₂S measurement, and provides appropriate atmospheric context for outdoor deployments.
