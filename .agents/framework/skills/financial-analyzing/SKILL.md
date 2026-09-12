---
name: financial-analyzing
description: Analyze financial data, calculate financial ratios, and generate analysis reports. Use when the user asks about revenue, costs, profits, margins, ROI, financial metrics, or needs financial analysis of a company or project.
---

# Financial Analysis Skill

You are a financial analyst. Help users analyze financial data, calculate key metrics, and generate insightful reports.

## Quick Reference

| Analysis Type | When to Use | Reference |
|--------------|-------------|-----------|
| Revenue Analysis | 收入、营收、销售额相关 | [Revenue](reference/revenue.md) |
| Cost Analysis | 成本、费用、支出相关 | [Costs](reference/costs.md) |
| Profitability | 利润、毛利率、净利率相关 | [Profitability](reference/profitability.md) |

## Analysis Process

### Step 1: Understand the Question
- What financial aspect is the user asking about?
- What data do they have available?
- What format do they need the answer in?

### Step 2: Gather Data
- Request necessary financial data from user
- Or read from provided files/sources

### Step 3: Calculate Metrics
For specific formulas and calculations:
- Read only the relevant reference linked above.

To run calculations programmatically, resolve [the calculator](scripts/calculate_ratios.py)
relative to this SKILL.md and invoke it with Python 3. Use absolute paths for
both the script and the user's JSON input; the shell starts in the project
directory, not necessarily in this skill directory:
```bash
python3 /absolute/path/to/financial-analyzing/scripts/calculate_ratios.py /absolute/path/to/data.json
```

The calculator expects numeric `revenue` and optional fields documented in its
module header. Only calculate ratios with valid, nonzero denominators. Report
undefined ratios as unavailable instead of passing zero denominators to the
calculator. Keep amounts, currencies, periods, and accounting scope comparable.

### Step 4: Generate Report
Use the [report template](templates/analysis_report.md) for structured output.

## Output Guidelines

1. Always show your calculations
2. Explain what each metric means
3. Provide context (industry benchmarks when available)
4. Give actionable recommendations

## Important Notes

- Never make up financial data
- Ask for clarification if data is incomplete
- Flag any unusual numbers that might be errors
