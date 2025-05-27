from sklearn.base import BaseEstimator, TransformerMixin


class PoiAndBoolPreprocessor(BaseEstimator, TransformerMixin):
    poi_list = [
        "school",
        "clinic",
        "postOffice",
        "kindergarten",
        "restaurant",
        "college",
        "pharmacy",
    ]
    bool_cols = [
        "hasParkingSpace",
        "hasBalcony",
        "hasElevator",
        "hasSecurity",
        "hasStorageRoom",
    ]
    drop_cols = ["latitude", "longitude", "ownership"]

    def fit(self, X, y=None):
        return self

    def transform(self, X):
        X = X.copy()
        X = X.drop(columns=self.drop_cols)

        for poi in self.poi_list:
            dist_col = poi + "Distance"
            flag_col = "has" + poi.capitalize()
            X[flag_col] = X[dist_col].apply(lambda x: 1 if x <= 1 else 0)
        X = X.drop(columns=[poi + "Distance" for poi in self.poi_list])

        for col in self.bool_cols:
            X[col] = X[col].apply(lambda x: 1 if x == "yes" else 0)

        X["age"] = 2025 - X["buildYear"]
        X = X.drop(columns=["buildYear"])

        return X
