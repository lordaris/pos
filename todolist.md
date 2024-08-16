TODO:

- Crear un campo "modified" para actualización de usuarios

- [ ] Crear una ruta para obtener un ticket por medio de su número de ticket.

- [ ] Considerar el número de caja para la impresión del ticket

- [ ] Considerar tener un almacenamiento de inicios de sesión

- [ ] Mandar desde los productos el descuento que se obtendría por una oferta (el buyget tiene una lógica especial).

- [ ] Ticket debe regresar el tipo de promoción por cada producto, precio regular y lo que corresponda segun el tipo de promoción:

  - Buyget:
    - Precio regular
    - Cantidad de productos pagados y su total
    - Cantidad de productos regalados y su total igual a cero
  - DiscountPercentage:
    - Precio regular
    - Porcentaje de descuento
    - Precio con descuento
    - Total con descuento
  - DiscountPrice:
    - Precio regular:
    - Precio de descuento:
    - Total con descuento.

- [ ] Considerar la eliminación de los permisos dentro de los roles y utilizar los puros roles para dar autorizaciones

- [ ] Actualizar el API para agregar administración financiera.

- [ ] Manejo de pedidos
