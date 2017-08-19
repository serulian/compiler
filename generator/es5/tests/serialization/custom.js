$module('custom', function () {
  var $static = this;
  this.$class('38a26367', 'CustomJSON', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function () {
      var instance = new $static();
      return instance;
    };
    $static.Get = function () {
      return $g.custom.CustomJSON.new();
    };
    $instance.Stringify = function (value) {
      var $this = this;
      return $g.________testlib.basictypes.JSON.Get().Stringify(value);
    };
    $instance.Parse = function (value) {
      var $this = this;
      return $g.________testlib.basictypes.JSON.Get().Parse(value);
    };
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Get|1|fd8bc7c9<38a26367>": true,
        "Stringify|2|fd8bc7c9<44e219a9>": true,
        "Parse|2|fd8bc7c9<ad6de9ce<any>>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('9a932bd6', 'AnotherStruct', false, '', function () {
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
        "Parse|1|fd8bc7c9<9a932bd6>": true,
        "equals|4|fd8bc7c9<54ff3ddf>": true,
        "Stringify|2|fd8bc7c9<44e219a9>": true,
        "Mapping|2|fd8bc7c9<ad6de9ce<any>>": true,
        "Clone|2|fd8bc7c9<9a932bd6>": true,
        "String|2|fd8bc7c9<44e219a9>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  this.$struct('f97a88eb', 'SomeStruct', false, '', function () {
    var $static = this;
    var $instance = this.prototype;
    $static.new = function (SomeField, AnotherField, SomeInstance) {
      var instance = new $static();
      instance[BOXED_DATA_PROPERTY] = {
        SomeField: SomeField,
        AnotherField: AnotherField,
        SomeInstance: SomeInstance,
      };
      instance.$markruntimecreated();
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
      return $g.custom.AnotherStruct;
    }, function () {
      return $g.custom.AnotherStruct;
    }, false);
    this.$typesig = function () {
      if (this.$cachedtypesig) {
        return this.$cachedtypesig;
      }
      var computed = {
        "Parse|1|fd8bc7c9<f97a88eb>": true,
        "equals|4|fd8bc7c9<54ff3ddf>": true,
        "Stringify|2|fd8bc7c9<44e219a9>": true,
        "Mapping|2|fd8bc7c9<ad6de9ce<any>>": true,
        "Clone|2|fd8bc7c9<f97a88eb>": true,
        "String|2|fd8bc7c9<44e219a9>": true,
      };
      return this.$cachedtypesig = computed;
    };
  });

  $static.TEST = $t.markpromising(function () {
    var $result;
    var jsonString;
    var parsed;
    var s;
    var $current = 0;
    var $continue = function ($resolve, $reject) {
      localasyncloop: while (true) {
        switch ($current) {
          case 0:
            s = $g.custom.SomeStruct.new($t.fastbox(2, $g.________testlib.basictypes.Integer), $t.fastbox(false, $g.________testlib.basictypes.Boolean), $g.custom.AnotherStruct.new($t.fastbox(true, $g.________testlib.basictypes.Boolean)));
            jsonString = $t.fastbox('{"AnotherField":false,"SomeField":2,"SomeInstance":{"AnotherBool":true}}', $g.________testlib.basictypes.String);
            $promise.maybe($g.custom.SomeStruct.Parse($g.custom.CustomJSON)(jsonString)).then(function ($result0) {
              $result = $result0;
              $current = 1;
              $continue($resolve, $reject);
              return;
            }).catch(function (err) {
              $reject(err);
              return;
            });
            return;

          case 1:
            parsed = $result;
            $resolve($t.fastbox(((parsed.SomeField.$wrapped == 2) && !parsed.AnotherField.$wrapped) && parsed.SomeInstance.AnotherBool.$wrapped, $g.________testlib.basictypes.Boolean));
            return;

          default:
            $resolve();
            return;
        }
      }
    };
    return $promise.new($continue);
  });
});
