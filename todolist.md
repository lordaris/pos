TODO:

- Continuar el invoice method, debe regresar todos los datos del ticket, y recibir a fuerzas los datos como
  el id del cajero. Además debe ir protegido con un middleware que permita solo el uso de roles de cajero, de gerente y de admin.
- Crear un campo "modified" para actualización de usuarios
- Validar que los tokens utilizados existan y estén activos. Si no, eliminarlos.
- Regresar el nombre de cada producto en el invoice.

- Modificar la estructura de tokens para almacenarla en la colección de usuarios
  - Evito tener tokens huerfanos y referencias rotas
  - Se maneja todo en la misma colección.
- Modificar los handlers relacionados con tokens y el middleware, ya que estos
  buscan los tokens en una colección separada.

- Crear un handler para los movimientos de inventario.

  - Ingresar un tipo de movimiento
  - Recibir el ID del producto
  - Buscar el barcode del producto
  - Utilizar el middleware y el context

- Modificar el total del ticket para que solo regrese dos decimales.

e pueda incluir el contexto y se registre el usuario (cajero) en el invoice
