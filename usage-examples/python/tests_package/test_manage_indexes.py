import unittest
import examples.indexes.create_basic
import examples.indexes.create_filter
import examples.indexes.view
import examples.indexes.drop
import time
from dotenv import load_dotenv

class TestManageIndexes(unittest.TestCase):
    def setUp(self):
        load_dotenv()

    def test_create_basic(self):
        print("----------Test examples should create a basic vector index----------")
        examples.indexes.create_basic.example()
        print("----------Test complete, cleaning up----------")
        examples.indexes.drop.example()

    def test_create_filter(self):
        print("----------Test examples should create a vector index with filter----------")
        examples.indexes.create_filter.example()
        print("----------Test complete, cleaning up----------")
        examples.indexes.drop.example()

    def test_view(self):
        print("----------Test examples should successfully view an index, setting up----------")
        examples.indexes.create_basic.example()
        time.sleep(5)
        print("----------Set up complete, running tests_package----------")
        indexes = examples.indexes.view.example()
        self.assertEqual(1, len(indexes))
        index_name = indexes[0].get("name")
        self.assertEqual("vector_index", index_name)
        latestIndexDefinition = indexes[0].get("latestDefinition")
        definitionFields = latestIndexDefinition.get("fields")
        indexPath = definitionFields[0].get("path")
        self.assertEqual("plot_embedding", indexPath)
        numDimensions = definitionFields[0].get("numDimensions")
        self.assertEqual(1536, numDimensions)
        similarity = definitionFields[0].get("similarity")
        self.assertEqual("euclidean", similarity)
        time.sleep(5)
        print("----------Test complete, cleaning up----------")
        examples.indexes.drop.example()

    def test_drop(self):
        print("----------Test examples should successfully drop an index, setting up----------")
        examples.indexes.create_basic.example()
        time.sleep(5)
        print("----------Set up complete, running tests_package----------")
        examples.indexes.drop.example()
