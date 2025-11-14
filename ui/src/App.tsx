import React, { useEffect, useState } from 'react'

type Row = { catalog: string, schema: string, table: string }

export default function App(){
  const [rows, setRows] = useState<Row[]>([])
  useEffect(() => {
    // very simple: show static guidance and example queries
    setRows([
      {catalog:'delta', schema:'default', table:'mobility-delta/bronze/trips'},
      {catalog:'gold', schema:'gold', table:'gold_market_hourly_kpis'}
    ])
  },[])

  return (
    <div style={{fontFamily:'system-ui', padding:16}}>
      <h1>Mobility Lakehouse — Data Catalog</h1>
      <p>Run these in Trino (http://localhost:8082):</p>
      <pre>
{`SHOW CATALOGS;
SHOW SCHEMAS FROM delta;
SELECT * FROM delta.default."mobility-delta/bronze/trips" LIMIT 5;
SELECT * FROM gold.gold_market_hourly_kpis ORDER BY hour DESC LIMIT 20;`}
      </pre>
      <h3>Known Tables</h3>
      <ul>
        {rows.map((r,i)=>(<li key={i}><code>{r.catalog}.{r.schema}."{r.table}"</code></li>))}
      </ul>
    </div>
  )
}
