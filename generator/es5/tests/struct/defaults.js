$module('defaults', function () {
  var $static = this;
  this.$struct('6cbd0ebf', 'AnotherStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (AnotherBool) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        AnotherBool: AnotherBool,
      };
      instance.$markruntimecreated();
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'AnotherBool', 'AnotherBool', function () {
      return $g.________testlib.basictypes.Boolean;
    }, function () {
      return $g.________testlib.basictypes.Boolean;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|fb1385bf<6cbd0ebf>": true,
        "equals|4|fb1385bf<71258460>": true,
        "Stringify|2|fb1385bf<b2b53db7>": true,
        "Mapping|2|fb1385bf<204295f9<any>>": true,
        "Clone|2|fb1385bf<6cbd0ebf>": true,
        "String|2|fb1385bf<b2b53db7>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('41f59c9b', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
      };
      instance.$markruntimecreated();
      return $static.$initDefaults(instance, true);
    };
    $static.$initDefaults = function (instance, isRuntimeCreated) {
      var boxed = instance[BOXED_DATA_PROPERTY];
      if (isRuntimeCreated || (boxed['SomeField'] === undefined)) {
        instance.SomeField = $t.fastbox(42, $g.________testlib.basictypes.Integer);
      }
      if (isRuntimeCreated || (boxed['AnotherField'] === undefined)) {
        instance.AnotherField = $t.fastbox(false, $g.________testlib.basictypes.Boolean);
      }
      if (isRuntimeCreated || (boxed['SomeInstance'] === undefined)) {
        instance.SomeInstance = $g.defaults.AnotherStruct.new($t.fastbox(true, $g.________testlib.basictypes.Boolean));
      }
      return instance;
    };
    $static.$fields = [];
    $t.defineStructField($static, 'SomeField', 'SomeField', function () {
      return $g.________testlib.basictypes.Integer;
    }, function () {
      return $g.________testlib.basictypes.Integer;
    }, false);
    $t.defineStructField($static, 'AnotherField', 'AnotherField', function () {
      return $g.________testlib.basictypes.Boolean;
    }, function () {
      return $g.________testlib.basictypes.Boolean;
    }, false);
    $t.defineStructField($static, 'SomeInstance', 'SomeInstance', function () {
      return $g.defaults.AnotherStruct;
    }, function () {
      return $g.defaults.AnotherStruct;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|fb1385bf<41f59c9b>": true,
        "equals|4|fb1385bf<71258460>": true,
        "Stringify|2|fb1385bf<b2b53db7>": true,
        "Mapping|2|fb1385bf<204295f9<any>>": true,
        "Clone|2|fb1385bf<41f59c9b>": true,
        "String|2|fb1385bf<b2b53db7>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = function () {
    var $temp0;
    var ss;
    ss = ($temp0 = $g.defaults.SomeStruct.new(), $temp0.AnotherField = $t.fastbox(true, $g.________testlib.basictypes.Boolean), $temp0);
    return $t.fastbox(((ss.SomeField.$wrapped == 42) && ss.AnotherField.$wrapped) && ss.SomeInstance.AnotherBool.$wrapped, $g.________testlib.basictypes.Boolean);
  };
});
