# TODO

- [ ] Implementar bien los validadores, ya hay uno empezado en data products

- Crear un campo "modified" para actualización de usuarios

- [ ] Actualizar el API para agregar administración financiera.

- [ ] Manejo de pedidos

- [ ] Considerar tener un almacenamiento de inicios de sesión y registro de acciones del usuario

- [ ] Hacer que el getUserByRole reciba el nombre del rol en lugar de el ID

- [ ] Considerar el número de caja para la impresión del ticket

## Caracteristicas

### Configuración

- [ ] Información del negocio
- [ ] ~~Facturación electronica~~
- [x] Roles
- [x] Usuarios
- [ ] Tienda

### Estadísticas

- [ ] Ventas
- [ ] Ventas por producto

### Contabilidad

- [ ] Gastos
- [ ] Tipos de gastos
- [ ] Impuestos
- [ ] Clientes

### Productos

- [x] Categorías
- [x] Productos
- [x] Promociones

### Inventario

- [x] Inventario
- [x] Movimientos de inventario
      ~~Historial de inventario~~

### Ventas

- [ ] Cuadre de caja

#### Ideas

1. Un API para la parte central
2. Puede obtener datos de todas las tiendas (nombre, ubicacion, Estadísticas, etc...)
3. Handler para obtener el inventario cada tienda
4. Handler para hacer movimientos de inventario

- debe recibir:
  - tienda de origen y de destino,
  - el array de productos y cantidades
  - el estado del movimiento.
- Las fechas se generan automaticamente.
  - La de inicio, se genera al crear el movimientos
  - la de finalizacion se genera cuando el estado de movimiento es == completado.
- El stock global y los locales se actualizan automaticamente
  - El local de donde se hace la transferencia se actualiza cuando el movimiento es en proceso
  - El local de destino se actualiza cuando el movimiento es completado

5. Handler estadísticas de ventas

- Recibir un periodo (fecha de inicio y de finalización)
- Filtrar las tiendas o usar todas

6. Handler para estadisticas de ventas por producto

- Recibir un periodo
- Filtrar tiendas
- Regresar el total de ventas por periodo
- Regresar el promedio de ventas por día
