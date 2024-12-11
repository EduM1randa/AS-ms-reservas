package controllers

import (
	"fmt"
	"strconv"
	"strings"

	db "gitlab.com/los-pelagatos-al-dia/pids-productos/database"
	"gitlab.com/los-pelagatos-al-dia/pids-productos/dto"
	"gitlab.com/los-pelagatos-al-dia/pids-productos/models"
	"gorm.io/gorm/clause"
)

func GetProducts(productsDto dto.ProductsDto) (models.Response, error) {

	name := productsDto.Name
	ids := productsDto.ProductsId
	brand := productsDto.Brand                                            // Ejemplo: ?brand=Algo
	category := productsDto.Category                                      // Ejemplo: ?category=Otra cosa
	tags := productsDto.Tags                                              // Ejemplo: ?tags=Tag1,Tag2,Tag3
	limit := ObtenerLimit(productsDto.Limit, DefaultProductLimit)         // Ejemplo: ?limit=20
	isEagle := productsDto.IsEagle            // Ejemplo: ?eagle=true
	minPriceParam := productsDto.MinPriceParam                            // Ejemplo: ?min_price=1000
	maxPriceParam := productsDto.MaxPriceParam                            // Ejemplo: ?max_price=5000
	includeImages := productsDto.IncludeImages // Ejemplo: ?include_images=true

	var products []models.Product // Productos que devolverá
	
	query := db.DB.Debug().Model(&models.Product{}) // Inicia la consulta

	if !includeImages {
		query = query.Omit("image")
	}

	if isEagle {
		query = query.Preload("Brand").Preload("Category").Preload("Tags")
	}

	if len(ids) > 0 {
		query = query.Where("products.id IN (?)", ids)
	}

	if brand != "" {
		query = query.Joins("JOIN brands ON products.brand_id = brands.id").
			Where("brands.name = ?", brand)
	}

	if category != "" {
		query = query.Joins("JOIN categories ON products.category_id = categories.id").
			Where("categories.name = ?", category)
	}

	if len(tags) > 0 {
		query = query.Joins("JOIN product_tags ON products.id = product_tags.product_id").
			Joins("JOIN tags ON product_tags.tag_id = tags.id").
			Where("tags.name IN (?)", tags)
	}

	if minPriceParam != "" {
		minPrice, err := strconv.Atoi(minPriceParam)
		if err != nil || minPrice <= 0 {

			return models.Response{
				Status:  "400",
				Success: "false",
				Message: "Precio mínimo inválido",
			}, fmt.Errorf("el precio mínimo debe ser mayor o igual a 0");

		}
		query = query.Where("price >= ?", minPrice)
	}

	if maxPriceParam != "" {
		maxPrice, err := strconv.Atoi(maxPriceParam)
		if err != nil || maxPrice <= 0 {

			return models.Response{
				Status:  "400",
				Success: "false",
				Message: "Precio máximo inválido",
			}, fmt.Errorf("el precio máximo debe ser mayor o igual a 0");

		}

		query = query.Where("price <= ?", maxPrice)

	}

	query = query.Limit(limit)

	if err := query.Find(&products).Error; err != nil {
		fmt.Println(err)

		return models.Response{
			Status:  "500",
			Success: "false",
			Message: "Error obteniendo productos",
		}, err;

	}

	if name != "" {
		productsFiltered := []models.Product{} // Productos que devolverá
		for i := 0; i < len(products); i++ {
			if strings.Contains(strings.ToLower(products[i].Name), strings.ToLower(name)) {
				productsFiltered = append(productsFiltered, products[i])
			}
		}

		return models.Response{
			Status:  "200",
			Success: "true",
			Message: "Success",
			Data:    productsFiltered,
		}, nil
	}

	return models.Response{
		Status:  "200",
		Success: "true",
		Message: "Success",
		Data:    products,
	}, nil
}

func GetProduct(productDto dto.ProductDto) (models.Response, error) {
	
	productId := productDto.ProductId
	isEagle := productDto.IsEagle         // Ejemplo: ?eagle=true
	includeImage := productDto.IncludeImage // Ejemplo: ?include_image=true

	var product models.Product
	query := db.DB.Model(&models.Product{})

	fmt.Println(includeImage)

	if !includeImage {
		query = query.Omit("image")
	}

	if isEagle {
		query = query.Preload("Brand").Preload("Category").Preload("Tags")
	}
	
	_ = query.Where("id = ?", productId).First(&product)

	if product.ID == 0 {
		return models.Response{
			Status:  "404",
			Success: "false",
			Message: "Producto no encontrado",
		}, fmt.Errorf("no se pudo encontrar un producto con esa id");
	}

	return models.Response{
		Status:  "200",
		Success: "true",
		Message: "Success",
		Data:    product,
	}, nil
}

func PostProduct(productDto dto.CreateProductDto) (models.Response, error) {
	fmt.Println("Entre a postear un producto");
	var product models.Product
	product.Name = productDto.Name
	product.Description = productDto.Description
	product.BrandID = productDto.BrandID
	product.CategoryID = productDto.CategoryID
	product.Price = productDto.Price
	
	_, errBrand := GetBrand(dto.BrandDto{BrandId: strconv.Itoa(productDto.BrandID)})
	if errBrand != nil {
		return models.Response{
			Status:  "404",
			Success: "false",
			Message: "Error creando producto, marca no existe",
		}, errBrand;
	}
	_, errCategory := GetCategory(dto.CategoryDto{CategoryId: strconv.Itoa(productDto.CategoryID)})
	if errCategory != nil {
		return models.Response{
			Status:  "404",
			Success: "false",
			Message: "Error creando producto, categoria no existe",
		}, errCategory;
	}

	db.DB.Create(&product)

	if product.ID == 0 {
		return models.Response{
			Status:  "404",
			Success: "false",
			Message: "Error creando producto",
		}, fmt.Errorf("no se pudo crear el producto");
	}

	return models.Response{
		Status:  "200",
		Success: "true",
		Message: "Success",
		Data:    product,
	}, nil
}

func UpdateProduct(productDto dto.UpdateProductDto) (models.Response, error) {
	fmt.Println("Entre a actualizar un producto");
	var updateProduct models.Product
	productId := strconv.Itoa(productDto.ProductId)
	updateProduct.Name = productDto.Name
	updateProduct.Description = productDto.Description
	updateProduct.BrandID = productDto.BrandID
	updateProduct.CategoryID = productDto.CategoryID
	updateProduct.Price = productDto.Price
	fmt.Println(productId);
	var product models.Product

	_ = db.DB.Where("id = ?", productId).First(&product).Updates(models.Product{Name: updateProduct.Name, Description: updateProduct.Description, 
		Price: updateProduct.Price, BrandID: updateProduct.BrandID, CategoryID: updateProduct.CategoryID});

	if product.ID == 0 {
		return models.Response{
			Status:  "404",
			Success: "false",
			Message: "Producto no encontrado",
		}, fmt.Errorf("no se encontró un producto con esa id");
	}

	return models.Response{
		Status:  "200",
		Success: "true",
		Message: "Success",
		Data:    product,
	}, nil
}

func UpdateProductTags(productDto dto.UpdateProductDto) (models.Response, error) {
	fmt.Println("Entre a actualizar un producto");
	var updateProduct models.Product
	productId := strconv.Itoa(productDto.ProductId)
	updateProduct.Name = productDto.Name
	fmt.Println(productId);
	var product models.Product

	_ = db.DB.Where("id = ?", productId).Association("Tags").Append(&models.Tag{});

	if product.ID == 0 {
		return models.Response{
			Status:  "404",
			Success: "false",
			Message: "Producto no encontrado",
		}, fmt.Errorf("no se encontró un producto con esa id");
	}

	return models.Response{
		Status:  "200",
		Success: "true",
		Message: "Success",
		Data:    product,
	}, nil
}

func DeleteProduct(productDto dto.ProductDto) (models.Response, error) {
	productId := productDto.ProductId

	var product models.Product
	db.DB.Clauses(clause.Returning{}).Where("id = ?", productId).Delete(&product)

	if product.ID == 0 {
		return models.Response{
			Status:  "404",
			Success: "false",
			Message: "Producto no encontrado",
		}, fmt.Errorf("no se encontró un producto con esa id");
	}

	return models.Response{
		Status:  "200",
		Success: "true",
		Message: "Success",
		Data:    product,
	}, nil
}

/*
INSERT INTO brand (name) VALUES
('Sanicat'),
('Royal Canin'),
('Pro Plan'),
('Kong'),
('Nath'),
('Fit Formula'),
('The Cat Band'),
('Play&Bite');

INSERT INTO category (name) VALUES
('Comida'),
('Accesorios'),
('Farmacia'),
('Juguetes');

INSERT INTO tag (name) VALUES
('perro'),
('gato'),
('pequeño animal'),
('alimento seco'),
('alimento humedo'),
('alimento medicado'),
('alimento necesidades especiales'),
('snacks'),
('arnes'),
('collar'),
('ropa'),
('limpieza'),
('medicamento'),
('arena de gato'),
('repelente'),
('shampoo'),
('juguete interactivo'),
('juguete funcional'),
('puzzle');

INSERT INTO products(name, description, brand_id, category_id, price) VALUES
('Nath adulto Maxi sabor pollo y arroz integral alimento para perros 3KG','Alimento para perros con sabor a pollo y arroz marrón ',4,6,19990),
('Royal Canin Alimento Húmedo Perro Adulto Recovery 145 G','Alimento en lata húmeda para perros y gatos convalecientes',1,6,3990),
('Arena para gatos Sanicat clumping 12 KG','Arena gruesa aglomerante para mantener las patas y los suelos limpios',8,8,17990);

INSERT INTO product_tag(product_id, tag_id) VALUES
(1,39),
(1,42),
(2,39),
(2,43),
(2,44),
(3,40),
(3,50),
(3,52);

SELECT * FROM category, brand;
SELECT * FROM tag;
SELECT * FROM product;

SELECT * FROM tags


	tagIDs := []uint{2, 12, 14}

	products := []*models.Product{
		{
			Name: "Arena para gatos Sanicat clumping 12 KG",
			Description: "Arena gruesa aglomerante para mantener las patas y los suelos limpios.",
			Price: 17990,
			BrandID: 4,
			CategoryID: 2,
		},
	}

	for _, tagID := range tagIDs {
		tag := &models.Tag{}
		tag.ID = tagID
		products[0].Tags = append(products[0].Tags, *tag)
	}

	resultProducts := db.DB.Create(products)
	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"message": "Success",
		"data":    resultProducts,
	})
*/
