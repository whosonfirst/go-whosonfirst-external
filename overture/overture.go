package overture

/*

D DESCRIBE( SELECT * FROM read_parquet('~/data/overture/parquet/*.parquet'));
┌─────────────┬────────────────────────────────────────────────────────────────────────────────────┬──────┬─────┬─────────┬───────┐
│ column_name │                                    column_type                                     │ null │ key │ default │ extra │
├─────────────┼────────────────────────────────────────────────────────────────────────────────────┼──────┼─────┼─────────┼───────┤
│ id          │ VARCHAR                                                                            │ YES  │     │         │       │
│ geometry    │ GEOMETRY                                                                           │ YES  │     │         │       │
│ bbox        │ STRUCT(xmin FLOAT, xmax FLOAT, ymin FLOAT, ymax FLOAT)                             │ YES  │     │         │       │
│ version     │ INTEGER                                                                            │ YES  │     │         │       │
│ sources     │ STRUCT(property VARCHAR, dataset VARCHAR, record_id VARCHAR, update_time VARCHA... │ YES  │     │         │       │
│ names       │ STRUCT("primary" VARCHAR, common MAP(VARCHAR, VARCHAR), rules STRUCT(variant VA... │ YES  │     │         │       │
│ categories  │ STRUCT("primary" VARCHAR, alternate VARCHAR[])                                     │ YES  │     │         │       │
│ confidence  │ DOUBLE                                                                             │ YES  │     │         │       │
│ websites    │ VARCHAR[]                                                                          │ YES  │     │         │       │
│ socials     │ VARCHAR[]                                                                          │ YES  │     │         │       │
│ emails      │ VARCHAR[]                                                                          │ YES  │     │         │       │
│ phones      │ VARCHAR[]                                                                          │ YES  │     │         │       │
│ brand       │ STRUCT(wikidata VARCHAR, "names" STRUCT("primary" VARCHAR, common MAP(VARCHAR, ... │ YES  │     │         │       │
│ addresses   │ STRUCT(freeform VARCHAR, locality VARCHAR, postcode VARCHAR, region VARCHAR, co... │ YES  │     │         │       │
└─────────────┴────────────────────────────────────────────────────────────────────────────────────┴──────┴─────┴─────────┴───────┘

*/
