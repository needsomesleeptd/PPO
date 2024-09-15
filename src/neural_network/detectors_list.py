from detection_scripts.schemes.wrong_arrows_directions_2 import *
from detection_scripts.schemes.wrong_if_subscr_2 import *
from detection_scripts.schemes.wrong_terminators import *

from detection_scripts.tables.checking_names import *
from detection_scripts.tables.table_err_detector import *

from detection_scripts.graphs.find_legend import * 
from detection_scripts.graphs.find_wrong_subs import *

from detection_scripts.formulas.checking_number_pos import * 
from detection_scripts.formulas.checking_number import * 
from detection_scripts.formulas.checking_syntax import * 

def create_schemes_detectors():
    arr_detector = ArrowsDestinationErrDetector()
    if_detector = IfSubscriptionErrDetector()
    term_detector = TerminatorsErrDetector()
    return [arr_detector,if_detector,term_detector]

def create_tables_detectors():
    table_name_detecotr = TableNameErrDetector()
    return [table_name_detecotr]

def create_graphs_detectors():
    legend_detector = LegendErrorDetector()
    subs_detector = AxisSubsErrorDetector()
    return [legend_detector,subs_detector]

def create_formulas_detectors():
    check_numer_pos = CheckingFormulasPositions()
    check_number_subs = CheckingFormulasSyntax()
    check_number_syntax = CheckingFormulasSyntax()
    return [check_number_subs,check_numer_pos,check_number_syntax]
def create_all_detectors():
    return create_schemes_detectors() + create_tables_detectors() + create_graphs_detectors() + create_formulas_detectors()



    

    


