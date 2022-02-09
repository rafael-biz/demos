﻿using System.Collections.Generic;
using System.Threading.Tasks;

namespace MyStore.Services.Products
{
    public interface IProductList
    {
        Task<IList<ProductDetailsDTO>> GetListAsync();
    }
}
