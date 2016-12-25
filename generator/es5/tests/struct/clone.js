$module('clone', function () {
  var $static = this;
  this.$struct('be48fdbe', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
        AnotherField: AnotherField,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return $g.____testlib.basictypes.Integer;
    }, function () {
      return $g.____testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherField', 'AnotherField', function () {
      return $g.____testlib.basictypes.Boolean;
    }, function () {
      return $g.____testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|29dc432d<be48fdbe>": true,
        "equals|4|29dc432d<43834c3f>": true,
        "Stringify|2|29dc432d<5cffd9b5>": true,
        "Mapping|2|29dc432d<df58fcbd<any>>": true,
        "Clone|2|29dc432d<be48fdbe>": true,
        "String|2|29dc432d<5cffd9b5>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $temp0;
    var first;
    var second;
    first = $g.clone.SomeStruct.new($t.fastbox(42, $g.____testlib.basictypes.Integer), $t.fastbox(false, $g.____testlib.basictypes.Boolean));
    second = ($temp0 = first.Clone(), $temp0.AnotherField = $t.fastbox(true, $g.____testlib.basictypes.Boolean), $temp0);
    return $t.fastbox(((second.SomeField.$wrapped == 42) && second.AnotherField.$wrapped) && !first.AnotherField.$wrapped, $g.____testlib.basictypes.Boolean);
  };
});
