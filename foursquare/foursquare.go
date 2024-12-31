package foursquare

/*

D DESCRIBE( SELECT * FROM read_parquet('~/data/foursquare/parquet/*.parquet'));
┌─────────────────────┬────────────────────────────────────────────────────────────┬──────┬─────┬─────────┬───────┐
│     column_name     │                        column_type                         │ null │ key │ default │ extra │
├─────────────────────┼────────────────────────────────────────────────────────────┼──────┼─────┼─────────┼───────┤
│ fsq_place_id        │ VARCHAR                                                    │ YES  │     │         │       │
│ name                │ VARCHAR                                                    │ YES  │     │         │       │
│ latitude            │ DOUBLE                                                     │ YES  │     │         │       │
│ longitude           │ DOUBLE                                                     │ YES  │     │         │       │
│ address             │ VARCHAR                                                    │ YES  │     │         │       │
│ locality            │ VARCHAR                                                    │ YES  │     │         │       │
│ region              │ VARCHAR                                                    │ YES  │     │         │       │
│ postcode            │ VARCHAR                                                    │ YES  │     │         │       │
│ admin_region        │ VARCHAR                                                    │ YES  │     │         │       │
│ post_town           │ VARCHAR                                                    │ YES  │     │         │       │
│ po_box              │ VARCHAR                                                    │ YES  │     │         │       │
│ country             │ VARCHAR                                                    │ YES  │     │         │       │
│ date_created        │ VARCHAR                                                    │ YES  │     │         │       │
│ date_refreshed      │ VARCHAR                                                    │ YES  │     │         │       │
│ date_closed         │ VARCHAR                                                    │ YES  │     │         │       │
│ tel                 │ VARCHAR                                                    │ YES  │     │         │       │
│ website             │ VARCHAR                                                    │ YES  │     │         │       │
│ email               │ VARCHAR                                                    │ YES  │     │         │       │
│ facebook_id         │ BIGINT                                                     │ YES  │     │         │       │
│ instagram           │ VARCHAR                                                    │ YES  │     │         │       │
│ twitter             │ VARCHAR                                                    │ YES  │     │         │       │
│ fsq_category_ids    │ VARCHAR[]                                                  │ YES  │     │         │       │
│ fsq_category_labels │ VARCHAR[]                                                  │ YES  │     │         │       │
│ geom                │ GEOMETRY                                                   │ YES  │     │         │       │
│ bbox                │ STRUCT(xmin DOUBLE, ymin DOUBLE, xmax DOUBLE, ymax DOUBLE) │ YES  │     │         │       │
└─────────────────────┴────────────────────────────────────────────────────────────┴──────┴─────┴─────────┴───────┘

*/
