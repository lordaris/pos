TODO:

- Continuar el invoice method,y recibir a fuerzas los datos como
  el id del cajero. Además debe ir protegido con un middleware que permita solo el uso de roles de cajero, de gerente y de admin.
- Crear un campo "modified" para actualización de usuarios

- Crear un handler para los movimientos de inventario.

  - Ingresar un tipo de movimiento IN / OUT
  - Recibir el ID del producto
  - Buscar el barcode del producto
  - Utilizar el middleware y el context
    - Para registrar al usuario

- [x] Modificar el total del ticket para que solo regrese dos decimales.

- [x] Invoice debe regresar los datos del ticket

- [x] Validar que los tokens utilizados existan y estén activos. Si no, eliminarlos.

- [x] Modificar la estructura de tokens para almacenarla en la colección de usuarios

- [x] Modificar los handlers relacionados con tokens y el middleware, ya que estos buscan los tokens en una colección separada.

-[x] Regresar el nombre de cada producto en el invoice.
