import unittest
import examples.indexes.create_basic
import examples.indexes.create_filter
import examples.indexes.drop
import examples.queries.ann_basic
import examples.queries.ann_filter
import time
import os
from dotenv import load_dotenv

class TestQueries(unittest.TestCase):
    def tearDown(self):
        print("----------Test complete, cleaning up----------")
        examples.indexes.drop.example()

    def setUp(self):
        load_dotenv()

    def test_ann_basic(self):
        ENV = os.getenv("ENV")
        print("----------Test example should perform an ANN query with expected results----------")
        print("----------Setting up vector index----------")
        examples.indexes.create_basic.example()
        time.sleep(5)
        print("----------Performing the query tests_package----------")
        result = examples.queries.ann_basic.example()
        if ENV=="Atlas":
            expectedOutput = [{'plot': 'A reporter, learning of time travelers visiting 20th century disasters, tries to change the history they know by averting upcoming disasters.', 'title': 'Thrill Seekers', 'score': 0.7892671227455139},
                {'plot': 'At the age of 21, Tim discovers he can travel in time and change what happens and has happened in his own life. His decision to make his world a better place by getting a girlfriend turns out not to be as easy as you might think.', 'title': 'About Time', 'score': 0.7843604683876038},
                {'plot': 'Hoping to alter the events of the past, a 19th century inventor instead travels 800,000 years into the future, where he finds humankind divided into two warring races.', 'title': 'The Time Machine', 'score': 0.7801066637039185},
                {'plot': "After using his mother's newly built time machine, Dolf gets stuck involuntary in the year 1212. He ends up in a children's crusade where he confronts his new friends with modern techniques...", 'title': 'Crusade in Jeans', 'score': 0.7789170742034912},
                {'plot': 'An officer for a security agency that regulates time travel, must fend for his life against a shady politician who has a tie to his past.', 'title': 'Timecop', 'score': 0.7771612405776978},
                {'plot': 'A time-travel experiment in which a robot probe is sent from the year 2073 to the year 1973 goes terribly wrong thrusting one of the project scientists, a man named Nicholas Sinclair into a...', 'title': 'A.P.E.X.', 'score': 0.7730885744094849},
                {'plot': "Agent J travels in time to M.I.B.'s early days in 1969 to stop an alien from assassinating his friend Agent K and changing history.", 'title': 'Men in Black 3', 'score': 0.7712380886077881},
                {'plot': 'Bound by a shared destiny, a teen bursting with scientific curiosity and a former boy-genius inventor embark on a mission to unearth the secrets of a place somewhere in time and space that exists in their collective memory.', 'title': 'Tomorrowland', 'score': 0.7669923901557922},
                {'plot': 'With the help of his uncle, a man travels to the future to try and bring his girlfriend back to life.', 'title': 'Love Story 2050', 'score': 0.7649372816085815},
                {'plot': 'A dimension-traveling wizard gets stuck in the 21st century because cell-phone radiation interferes with his magic. With his home world on the brink of war, he seeks help from a jaded ...', 'title': 'The Portal', 'score': 0.7640786170959473}]
        elif ENV=="local":
            expectedOutput = [{'plot': 'A reporter, learning of time travelers visiting 20th century disasters, tries to change the history they know by averting upcoming disasters.', 'title': 'Thrill Seekers', 'score': 0.7892671227455139},
            {'plot': 'At the age of 21, Tim discovers he can travel in time and change what happens and has happened in his own life. His decision to make his world a better place by getting a girlfriend turns out not to be as easy as you might think.', 'title': 'About Time', 'score': 0.7843604683876038},
            {'plot': 'Hoping to alter the events of the past, a 19th century inventor instead travels 800,000 years into the future, where he finds humankind divided into two warring races.', 'title': 'The Time Machine', 'score': 0.7801067233085632},
            {'plot': "After using his mother's newly built time machine, Dolf gets stuck involuntary in the year 1212. He ends up in a children's crusade where he confronts his new friends with modern techniques...", 'title': 'Crusade in Jeans', 'score': 0.7789170742034912},
            {'plot': 'An officer for a security agency that regulates time travel, must fend for his life against a shady politician who has a tie to his past.', 'title': 'Timecop', 'score': 0.7771613597869873},
            {'plot': 'A time-travel experiment in which a robot probe is sent from the year 2073 to the year 1973 goes terribly wrong thrusting one of the project scientists, a man named Nicholas Sinclair into a...', 'title': 'A.P.E.X.', 'score': 0.7730885744094849},
            {'plot': "Agent J travels in time to M.I.B.'s early days in 1969 to stop an alien from assassinating his friend Agent K and changing history.", 'title': 'Men in Black 3', 'score': 0.7712380290031433},
            {'plot': 'Bound by a shared destiny, a teen bursting with scientific curiosity and a former boy-genius inventor embark on a mission to unearth the secrets of a place somewhere in time and space that exists in their collective memory.', 'title': 'Tomorrowland', 'score': 0.7669923901557922},
            {'plot': 'With the help of his uncle, a man travels to the future to try and bring his girlfriend back to life.', 'title': 'Love Story 2050', 'score': 0.7649372220039368},
            {'plot': 'A dimension-traveling wizard gets stuck in the 21st century because cell-phone radiation interferes with his magic. With his home world on the brink of war, he seeks help from a jaded ...', 'title': 'The Portal', 'score': 0.7640786170959473}]
        self.assertEqual(expectedOutput, result)

    def test_ann_filter(self):
        ENV = os.getenv("ENV")
        print("----------Test example should perform an ANN query with a filter with expected results----------")
        print("----------Setting up vector index----------")
        examples.indexes.create_filter.example()
        time.sleep(5)
        print("----------Performing the query tests_package----------")
        result = examples.queries.ann_filter.example()
        if ENV=="Atlas":
            expectedOutput = [{'plot': 'In this magical tale about the boy who refuses to grow up, Peter Pan and his mischievous fairy sidekick Tinkerbell visit the nursery of Wendy, Michael, and John Darling. With a sprinkling ...', 'title': 'Peter Pan', 'year': 1960, 'score': 0.748110830783844},
                {'plot': 'A down-on-his-luck inventor turns a broken-down Grand Prix car into a fancy vehicle for his children, and then they go off on a magical fantasy adventure to save their grandfather in a far-off land.', 'title': 'Chitty Chitty Bang Bang', 'year': 1968, 'score': 0.7442465424537659},
                {'plot': 'A young man comes to the rescue of his girlfriend abducted by thieves and brought to Rio. An extravagant adventure ensues.', 'title': 'That Man from Rio', 'year': 1964, 'score': 0.7416020035743713},
                {'plot': 'A boy raised by wolves tries to adapt to human village life.', 'title': 'Jungle Book', 'year': 1942, 'score': 0.7387760877609253},
                {'plot': 'A pilot, stranded in the desert, meets a little boy who is a prince on a planet.', 'title': 'The Little Prince', 'year': 1974, 'score': 0.7378944158554077},
                {'plot': 'A red balloon with a life of its own follows a little boy around the streets of Paris.', 'title': 'The Red Balloon', 'year': 1956, 'score': 0.7342712879180908},
                {'plot': 'A poor boy wins the opportunity to tour the most eccentric and wonderful candy factory of all.', 'title': 'Willy Wonka & the Chocolate Factory', 'year': 1971, 'score': 0.7342107892036438},
                {'plot': 'An apprentice witch, three kids and a cynical conman search for the missing component to a magic spell useful to the defense of Britain.', 'title': 'Bedknobs and Broomsticks', 'year': 1971, 'score': 0.7339356541633606},
                {'plot': 'Arriving home to find his native land under the yoke of corrupt merchants, an adventurer named Sadko sets sail in search of a mythical bird of happiness.', 'title': 'Sadko', 'year': 1953, 'score': 0.7339220643043518},
                {'plot': "A young boys' coming of age tale set in a strange, carnivalesque village becomes the recreation of a memory that the director has twenty years later.", 'title': 'Pastoral Hide and Seek', 'year': 1974, 'score': 0.733299970626831}]
        elif ENV=="local":
            expectedOutput = [{'plot': 'In this magical tale about the boy who refuses to grow up, Peter Pan and his mischievous fairy sidekick Tinkerbell visit the nursery of Wendy, Michael, and John Darling. With a sprinkling ...', 'title': 'Peter Pan', 'year': 1960, 'score': 0.748110830783844},
                {'plot': 'A down-on-his-luck inventor turns a broken-down Grand Prix car into a fancy vehicle for his children, and then they go off on a magical fantasy adventure to save their grandfather in a far-off land.', 'title': 'Chitty Chitty Bang Bang', 'year': 1968, 'score': 0.7442465424537659},
                {'plot': 'A young man comes to the rescue of his girlfriend abducted by thieves and brought to Rio. An extravagant adventure ensues.', 'title': 'That Man from Rio', 'year': 1964, 'score': 0.7416019439697266},
                {'plot': 'A boy raised by wolves tries to adapt to human village life.', 'title': 'Jungle Book', 'year': 1942, 'score': 0.7387760281562805},
                {'plot': 'A pilot, stranded in the desert, meets a little boy who is a prince on a planet.', 'title': 'The Little Prince', 'year': 1974, 'score': 0.7378944158554077},
                {'plot': 'A red balloon with a life of its own follows a little boy around the streets of Paris.', 'title': 'The Red Balloon', 'year': 1956, 'score': 0.7342712879180908},
                {'plot': 'A poor boy wins the opportunity to tour the most eccentric and wonderful candy factory of all.', 'title': 'Willy Wonka & the Chocolate Factory', 'year': 1971, 'score': 0.7342106699943542},
                {'plot': 'An apprentice witch, three kids and a cynical conman search for the missing component to a magic spell useful to the defense of Britain.', 'title': 'Bedknobs and Broomsticks', 'year': 1971, 'score': 0.7339357137680054},
                {'plot': 'Arriving home to find his native land under the yoke of corrupt merchants, an adventurer named Sadko sets sail in search of a mythical bird of happiness.', 'title': 'Sadko', 'year': 1953, 'score': 0.7339221239089966},
                {'plot': "A young boys' coming of age tale set in a strange, carnivalesque village becomes the recreation of a memory that the director has twenty years later.", 'title': 'Pastoral Hide and Seek', 'year': 1974, 'score': 0.7332999110221863}]
        self.assertEqual(expectedOutput, result)
