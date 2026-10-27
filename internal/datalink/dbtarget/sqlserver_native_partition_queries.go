package dbtarget

// Partition identities and owner values are bound parameters. Database-wide
// metadata visibility is required before a zero reference count proves isolation.
const sqlServerNativeObjectsQuery = `/* gw_native_partition_objects */
SELECT pf.function_id, pf.boundary_value_on_right, pf.type, ty.name, pp.scale,
       (SELECT COUNT(*) FROM sys.partition_parameters p WHERE p.function_id=pf.function_id),
       ps.data_space_id, ps.function_id,
       CONVERT(nvarchar(max), fowner.value), CONVERT(nvarchar(max), sowner.value),
       (SELECT COUNT(*) FROM sys.partition_schemes s WHERE s.function_id=pf.function_id),
       CONVERT(bit, HAS_PERMS_BY_NAME(DB_NAME(), 'DATABASE', 'VIEW DEFINITION'))
FROM (VALUES (1)) AS seed(n)
LEFT JOIN sys.partition_functions pf ON pf.name=@p1
LEFT JOIN sys.partition_parameters pp ON pp.function_id=pf.function_id AND pp.parameter_id=1
LEFT JOIN sys.types ty ON ty.user_type_id=pp.user_type_id AND ty.is_user_defined=0
LEFT JOIN sys.partition_schemes ps ON ps.name=@p2
LEFT JOIN sys.extended_properties fowner ON fowner.class=21 AND fowner.major_id=pf.function_id
    AND fowner.minor_id=0 AND fowner.name=@p3
LEFT JOIN sys.extended_properties sowner ON sowner.class=20 AND sowner.major_id=ps.data_space_id
    AND sowner.minor_id=0 AND sowner.name=@p3`

const sqlServerNativeBoundsQuery = `/* gw_native_partition_bounds */
SELECT boundary_id, parameter_id, CONVERT(datetime2(6), value)
FROM sys.partition_range_values WHERE function_id=@p1 ORDER BY boundary_id`

const sqlServerNativeDestinationsQuery = `/* gw_native_partition_destinations */
SELECT d.destination_id, ds.name, ds.type
FROM sys.destination_data_spaces d
JOIN sys.data_spaces ds ON ds.data_space_id=d.data_space_id
WHERE d.partition_scheme_id=@p1 ORDER BY d.destination_id`

const sqlServerNativeIndexesQuery = `/* gw_native_partition_indexes */
SELECT i.object_id, i.index_id, i.data_space_id,
       SUM(CASE WHEN ic.partition_ordinal=1 AND c.name=N'bucket_start' THEN 1 ELSE 0 END),
       SUM(CASE WHEN ic.key_ordinal>0 AND c.name=N'bucket_start' THEN 1 ELSE 0 END),
       i.is_unique
FROM sys.indexes i
LEFT JOIN sys.index_columns ic ON ic.object_id=i.object_id AND ic.index_id=i.index_id
LEFT JOIN sys.columns c ON c.object_id=ic.object_id AND c.column_id=ic.column_id
WHERE i.object_id=@p1 OR i.data_space_id=@p2
GROUP BY i.object_id, i.index_id, i.data_space_id, i.is_unique
ORDER BY i.object_id, i.index_id`
